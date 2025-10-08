package main

import "fmt"

// TODO: Go through the whole code and wrap context around errors instead of just returning them raw.

func NewAction(hyprctl hyprctl, dispatcher dispatcher) *Action {
	return &Action{
		hyprctl:    hyprctl,
		dispatcher: dispatcher,
	}
}

// When {1, 3} and we call GoTo index 2, it's creating a new workspace named 2 in between 1 and 3, instead of going to 3.
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

	// TODO: Make this compact flag configurable and optional (enabled by default).
	compact := true
	if compact {
		targetWsName, err := GetZeroWidthNameFromIndex(monitorID, targetWsIndex)

		if err != nil {
			return err
		}

		err = CompactLocalWorkspacesOnMonitor(a, monitorID, false)

		if err != nil {
			return err
		}

		return dispatcher.GoToWorkspace(targetWsName)
	}

	return dispatcher.GoToWorkspace(sortedLocalWs[targetWsIndex].Name)
}

func (a *Action) MoveToWorkspace(targetIndex int, all bool) error {
	hyprctl, dispatcher := a.hyprctl, a.dispatcher

	activeWs, err := hyprctl.GetActiveWorkspace()

	if err != nil {
		return err
	}

	monitorID := activeWs.MonitorID

	// Compact empty workspaces here?
	if false {
		err := CompactLocalWorkspacesOnMonitor(a, monitorID, false)

		if err != nil {
			return err
		}
	}

	sortedLocalWs, err := GetSortedWorkspacesOnMonitor(hyprctl, monitorID)

	if err != nil {
		return err
	}

	currentWsIndex := GetWorkspaceIndexOnList(sortedLocalWs, activeWs.ID)

	if currentWsIndex == -1 {
		return fmt.Errorf("current workspace (ID %d) not found in local workspace list", activeWs.ID)
	}

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
	monitors, err := a.hyprctl.GetMonitors()

	if err != nil {
		return err
	}

	for _, mon := range monitors {
		err := CompactLocalWorkspacesOnMonitor(a, mon.ID, true)

		if err != nil {
			return err
		}
	}

	return nil
}
