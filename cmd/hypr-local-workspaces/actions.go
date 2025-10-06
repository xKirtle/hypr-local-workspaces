package main

import "fmt"

func NewAction(hyprctl hyprctl, dispatcher dispatcher) *Action {
	return &Action{
		hyprctl:    hyprctl,
		dispatcher: dispatcher,
	}
}

func (a *Action) GoToWorkspace(targetIndex int) error {
	hyprctl, dispatcher := a.hyprctl, a.dispatcher

	activeWs, err := hyprctl.GetActiveWorkspace()

	if err != nil {
		return err
	}

	monitorID := activeWs.MonitorID
	sortedLocalWs, err := GetSortedWorkspacesOnMonitor(hyprctl, monitorID)

	if err != nil {
		return err
	}

	currentWsIndex := GetWorkspaceIndexOnList(sortedLocalWs, activeWs.ID)

	if currentWsIndex == -1 {
		return fmt.Errorf("current workspace (ID %d) not found in local workspace list", activeWs.ID)
	}

	// Compact empty workspaces here?

	targetWsIndex, _ := DecideTargetWorkspaceIndex(currentWsIndex, targetIndex, sortedLocalWs)

	if currentWsIndex == targetWsIndex {
		// No-op
		return nil
	}

	targetWsName, err := GetZeroWidthNameFromIndex(monitorID, targetWsIndex)

	if err != nil {
		return err
	}

	return dispatcher.GoToWorkspace(targetWsName)
}

func (a *Action) MoveToWorkspace(targetIndex int, all bool) error {
	hyprctl, dispatcher := a.hyprctl, a.dispatcher

	activeWs, err := hyprctl.GetActiveWorkspace()

	if err != nil {
		return err
	}

	monitorID := activeWs.MonitorID
	sortedLocalWs, err := GetSortedWorkspacesOnMonitor(hyprctl, monitorID)

	if err != nil {
		return err
	}

	currentWsIndex := GetWorkspaceIndexOnList(sortedLocalWs, activeWs.ID)

	if currentWsIndex == -1 {
		return fmt.Errorf("current workspace (ID %d) not found in local workspace list", activeWs.ID)
	}

	// Compact empty workspaces here?

	targetWsIndex, _ := DecideTargetWorkspaceIndex(currentWsIndex, targetIndex, sortedLocalWs)

	if currentWsIndex == targetWsIndex {
		// No-op
		return nil
	}

	targetWsName, err := GetZeroWidthNameFromIndex(monitorID, targetWsIndex)

	if err != nil {
		return err
	}

	if all {
		return dispatcher.MoveAllToWorkspace(targetWsName)
	} else {
		activeWindow, err := hyprctl.GetActiveWindow()

		if err != nil {
			return err
		}

		return dispatcher.MoveToWorkspace(targetWsName, activeWindow.Address)
	}
}

func (a *Action) CycleWorkspace(direction string) error {
	hyprctl, dispatcher := a.hyprctl, a.dispatcher

	activeWs, err := hyprctl.GetActiveWorkspace()

	if err != nil {
		return err
	}

	monitorID := activeWs.MonitorID
	sortedLocalWs, err := GetSortedWorkspacesOnMonitor(hyprctl, monitorID)

	if err != nil {
		return err
	}

	currentWsIndex := GetWorkspaceIndexOnList(sortedLocalWs, activeWs.ID)

	if currentWsIndex == -1 {
		return fmt.Errorf("current workspace (ID %d) not found in local workspace list", activeWs.ID)
	}

	// Compact empty workspaces here?

	dir := 1
	if direction == "prev" {
		dir = -1
	}

	targetWsIndex, _ := DecideTargetWorkspaceIndex(currentWsIndex, currentWsIndex+dir, sortedLocalWs)

	if currentWsIndex == targetWsIndex {
		// No-op
		return nil
	}

	targetWsName, err := GetZeroWidthNameFromIndex(monitorID, targetWsIndex)

	if err != nil {
		return err
	}

	return dispatcher.GoToWorkspace(targetWsName)

}

func (a *Action) InitWorkspaces() error {
	return nil
}
