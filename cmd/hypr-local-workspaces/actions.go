package main

import "fmt"

// TODO: Go through the whole code and wrap context around errors instead of just returning them raw.

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

	if all && activeWs.WindowsCount > 1 {
		clients, err := hyprctl.GetClientsInWorkspace(activeWs.ID)

		if err != nil {
			return err
		}

		for _, client := range clients {
			err := dispatcher.MoveAddrToWorkspace(targetWsName, client.Address)

			if err != nil {
				return err
			}
		}
	} else {
		err = dispatcher.MoveToWorkspace(targetWsName)

		if err != nil {
			return err
		}
	}

	// TODO: Make this compact flag configurable and optional (enabled by default).
	compact := true
	if compact && (activeWs.WindowsCount == 1 || all) {
		err := CompactLocalWorkspacesOnMonitor(a, monitorID, false)

		if err != nil {
			return err
		}
	}

	return nil
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
