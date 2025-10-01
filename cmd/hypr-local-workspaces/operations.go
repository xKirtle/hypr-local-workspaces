package main

import (
	"fmt"
	"sort"
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

// IMPROVEMENT: Workspace ordering relies on either parsed slot numbers or Hyprland IDs;
// numbering may differ from the user's desired custom ordering in niche setups.

func InitWorkspaces() error {
	monitors, err := fetchMonitors()
	if err != nil {
		return fmt.Errorf("error fetching monitors: %w", err)
	}

	if len(monitors) == 0 {
		return nil
	}

	workspaces, err := fetchWorkspaces()
	if err != nil {
		return fmt.Errorf("error fetching workspaces: %w", err)
	}

	// remember focused monitor to restore later
	focusedMonitorID := monitors[0].ID

	for _, monitor := range monitors {
		if monitor.Focused {
			focusedMonitorID = monitor.ID
			break
		}
	}

	for _, monitor := range monitors {
		localWorkspaces, err := GetSortedLocalWorkspaces(workspaces, monitor.ID)
		if err != nil {
			localWorkspaces = make([]WorkspaceDTO, 0, len(workspaces))
			for _, workspace := range workspaces {
				if workspace.MonitorID == monitor.ID {
					localWorkspaces = append(localWorkspaces, workspace)
				}
			}

			sort.Slice(localWorkspaces, func(i, j int) bool {
				return localWorkspaces[i].ID < localWorkspaces[j].ID
			})
		}

		err = CompactLocalWorkspacesSimple(monitor, localWorkspaces)
		if err != nil {
			return fmt.Errorf("error compacting workspaces on monitor %d: %w", monitor.ID, err)
		}
	}

	return HyprctlFocusMonitor(focusedMonitorID)
}
