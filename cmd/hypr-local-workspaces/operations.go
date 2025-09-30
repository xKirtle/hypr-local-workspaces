package main

import (
	"fmt"
)

// TODO: Refactor to reduce code duplication between GoToWorkspace, MoveToWorkspace and CycleWorkspace

func GoToWorkspace(targetWorkspace int) error {
	snapshot, err := TakeSnapshot(MaskGoto)
	if err != nil {
		return fmt.Errorf("error taking snapshot: %w", err)
	}

	monitor, ok := GetFocusedMonitor(snapshot.Monitors, snapshot.ActiveWorkspace)
	if !ok {
		return fmt.Errorf("no focused monitor found")
	}

	workspaces, err := GetSortedLocalWorkspaces(snapshot.Workspaces, monitor.ID)
	if err != nil {
		return fmt.Errorf("error getting sorted local workspaces: %w", err)
	}

	activeLocalIndex, err := ActiveLocalIndex(workspaces, snapshot.ActiveWorkspace)
	if err != nil {
		return fmt.Errorf("error getting active local index: %w", err)
	}

	targetIndex, noop := DecideGoToTargetIndex(targetWorkspace, workspaces, activeLocalIndex)
	if noop {
		fmt.Println(noop)
		return nil
	}

	err = CompactLocalWorkspacesSimple(monitor, workspaces)
	if err != nil {
		return fmt.Errorf("error compacting local workspaces: %w", err)
	}

	workspaceName := TargetNameForWorkspace(monitor.ID, targetIndex+1)
	return HyprctlWorkspace(workspaceName)
}

func MoveToWorkspace(targetWorkspace int, moveAll bool) error {
	snapshot, err := TakeSnapshot(MaskMove)
	if err != nil {
		return fmt.Errorf("error taking snapshot: %w", err)
	}

	monitor, ok := GetFocusedMonitor(snapshot.Monitors, snapshot.ActiveWorkspace)
	if !ok {
		return fmt.Errorf("no focused monitor found")
	}

	workspaces, err := GetSortedLocalWorkspaces(snapshot.Workspaces, monitor.ID)
	if err != nil {
		return fmt.Errorf("error getting sorted local workspaces: %w", err)
	}

	activeLocalIndex, err := ActiveLocalIndex(workspaces, snapshot.ActiveWorkspace)
	if err != nil {
		return fmt.Errorf("error getting active local index: %w", err)
	}

	targetIndex, noop := DecideGoToTargetIndex(targetWorkspace, workspaces, activeLocalIndex)
	if noop {
		fmt.Println(noop)
		return nil
	}

	targetWorkspaceName := TargetNameForWorkspace(monitor.ID, targetIndex+1)
	if !moveAll {
		err = HyprctlMoveToWorkspace(targetWorkspaceName, snapshot.ActiveWindow.Address)
	} else {
		clientsOnWorkspace := GetClientsOnWorkspace(snapshot.ActiveWorkspace.ID, snapshot.Clients)
		err = HyprctlMoveToWorkspaceAll(targetWorkspaceName, clientsOnWorkspace)
	}

	if err != nil {
		return fmt.Errorf("error moving window(s) to workspace: %w", err)
	}

	// Refresh workspaces after move so we can compact correctly
	// IMPROVEMENT: This refresh + compact is awfully inefficient
	refreshedSnapshot, err := TakeSnapshot(FWorkspaces)
	if err != nil {
		return fmt.Errorf("error taking workspace snapshot: %w", err)
	}

	workspaces, err = GetSortedLocalWorkspaces(refreshedSnapshot.Workspaces, monitor.ID)
	if err != nil {
		return fmt.Errorf("error getting sorted local workspaces: %w", err)
	}

	err = CompactLocalWorkspacesSimple(monitor, workspaces)
	if err != nil {
		return fmt.Errorf("error compacting local workspaces: %w", err)
	}

	return nil
}

func CycleWorkspace(direction string) error {
	snapshot, err := TakeSnapshot(MaskCycle)
	if err != nil {
		return fmt.Errorf("error taking snapshot: %w", err)
	}

	monitor, ok := GetFocusedMonitor(snapshot.Monitors, snapshot.ActiveWorkspace)
	if !ok {
		return fmt.Errorf("no focused monitor found")
	}

	workspaces, err := GetSortedLocalWorkspaces(snapshot.Workspaces, monitor.ID)
	if err != nil {
		return fmt.Errorf("error getting sorted local workspaces: %w", err)
	}

	activeLocalIndex, err := ActiveLocalIndex(workspaces, snapshot.ActiveWorkspace)
	if err != nil {
		return fmt.Errorf("error getting active local index: %w", err)
	}

	var targetIndex int
	if direction == "up" {
		targetIndex = activeLocalIndex + 1
	} else {
		targetIndex = activeLocalIndex - 1
	}

	targetIndex, noop := DecideGoToTargetIndex(targetIndex+1, workspaces, activeLocalIndex)
	if noop {
		fmt.Println(noop)
		return nil
	}

	err = CompactLocalWorkspacesSimple(monitor, workspaces)
	if err != nil {
		return fmt.Errorf("error compacting local workspaces: %w", err)
	}

	workspaceName := TargetNameForWorkspace(monitor.ID, targetIndex+1)
	return HyprctlWorkspace(workspaceName)
}

// IMPROVEMENT: Not very robust. Assumes default workspace behavior of Hyprland
// (one workspace per monitor, workspaces named "1", "2", ... on each monitor)
// If there are custom-named workspaces or multiple workspaces per monitor, this may misbehave

func InitWorkspaces() error {
	monitors, err := fetchMonitors()
	if err != nil {
		return fmt.Errorf("error fetching monitors: %w", err)
	}

	for _, monitor := range monitors {
		err := HyprctlFocusMonitor(monitor.ID)
		if err != nil {
			return fmt.Errorf("error focusing monitor %d: %w", monitor.ID, err)
		}

		workspace, err := fetchActiveWorkspace()
		if err != nil {
			return fmt.Errorf("error fetching active workspace on monitor %d: %w", monitor.ID, err)
		}

		targetName := TargetNameForWorkspace(monitor.ID, 1)

		if workspace.Name != targetName {
			err = HyprctlRenameWorkspace(workspace.ID, targetName)
			if err != nil {
				return fmt.Errorf("error renaming workspace %q to %q: %w", workspace.Name, targetName, err)
			}
		}
	}

	return HyprctlFocusMonitor(monitors[0].ID)
}
