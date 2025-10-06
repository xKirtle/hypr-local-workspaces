package main

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
		return nil
	}

	// Compact empty workspaces here?

	targetWsIndex := DecideTargetWorkspaceIndex(currentWsIndex, targetIndex, sortedLocalWs)

	if targetWsIndex == -1 {
		return nil
	}

	targetWsName, err := GetZeroWidthNameFromIndex(monitorID, targetWsIndex)

	if err != nil {
		return err
	}

	return dispatcher.GoToWorkspace(targetWsName)
}

func (a *Action) MoveToWorkspace(targetIndex int, all bool) error {
	return nil
}

func (a *Action) CycleWorkspace(dir string) error {
	return nil
}

func (a *Action) InitWorkspaces() error {
	return nil
}
