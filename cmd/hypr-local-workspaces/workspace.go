package main

import (
	"fmt"
	"sort"
	"strconv"
)

func GetFocusedMonitor(monitors []MonitorDTO, activeWorkspace WorkspaceDTO) (MonitorDTO, bool) {
	if len(monitors) == 1 {
		return monitors[0], true
	}

	// Trust Hyprland's focused flag
	for _, monitor := range monitors {
		if monitor.Focused {
			return monitor, true
		}
	}

	// If, for some reason, Hyprland focused flag is not being updated or doesn't exist,
	// return the monitor that hosts the active workspace.
	if activeWorkspace.Monitor != "" {
		for _, monitor := range monitors {
			if monitor.Name == activeWorkspace.Monitor {
				return monitor, true
			}
		}
	}

	return MonitorDTO{}, false
}

func GetSortedLocalWorkspaces(workspaces []WorkspaceDTO, monitorID int) ([]WorkspaceDTO, error) {
	type tmpWorkspaceDTO struct {
		workspace      WorkspaceDTO
		workspaceIndex int
	}

	tmp := make([]tmpWorkspaceDTO, 0, len(workspaces))
	for _, workspace := range workspaces {
		if workspace.MonitorID != monitorID {
			continue
		}
		workspaceIndex, err := ParseLocalWorkspace(workspace.Name)
		if err != nil {
			return nil, err
		}

		tmp = append(tmp, tmpWorkspaceDTO{workspace, workspaceIndex})
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].workspaceIndex < tmp[j].workspaceIndex
	})

	result := make([]WorkspaceDTO, len(tmp))
	for i := range tmp {
		result[i] = tmp[i].workspace
	}

	return result, nil
}

func ActiveLocalIndex(localWorkspaces []WorkspaceDTO, active WorkspaceDTO) (int, error) {
	activeSlot, err := ParseLocalWorkspace(active.Name)
	if err != nil {
		return -1, err
	}

	for i, w := range localWorkspaces {
		s, err := ParseLocalWorkspace(w.Name)
		if err != nil {
			return -1, err
		}
		if s == activeSlot {
			return i, nil
		}
	}

	return -1, nil
}

func LastOccupiedLocalIndex(localWorkspaces []WorkspaceDTO) int {
	last := -1
	for i, w := range localWorkspaces {
		if w.Windows > 0 {
			last = i
		}
	}

	return last
}

// DecideGoToTargetIndex returns (targetIndex, noOp).
// 0-based indexes; creation is signaled by targetIndex == len(localWorkspaces).
func DecideGoToTargetIndex(requested int, localWorkspaces []WorkspaceDTO, curIndex int) (int, bool) {
	requestedIndex := requested - 1 // convert to 0-based

	if requestedIndex == curIndex {
		return requested, true
	}

	if requestedIndex < 0 {
		return 0, true
	}

	lastOcc := LastOccupiedLocalIndex(localWorkspaces) // -1 if none
	boundary := lastOcc + 1

	// We allow at most:
	// - focus up to boundary (which may be an existing empty index)
	// - and, if boundary == len(localWorkspaces), allow creation at exactly len(localWorkspaces)
	if boundary > len(localWorkspaces) { // should only be == or <, but safe
		boundary = len(localWorkspaces)
	}

	// Clamp
	target := requestedIndex
	if target > boundary {
		target = boundary
	}

	if target < 0 {
		target = 0
	}

	// empty-upward guard: don't move up from an empty current
	if curIndex >= 0 && target > curIndex && localWorkspaces[curIndex].Windows == 0 {
		return curIndex, true
	}

	// same index -> no-op
	if curIndex >= 0 && target == curIndex {
		return target, true
	}

	return target, false
}

// CompactLocalWorkspacesSimple renames local workspaces on the given monitor to be sequentially numbered from 1.
func CompactLocalWorkspacesSimple(monitorDTO MonitorDTO, localWorkspaces []WorkspaceDTO) error {
	if len(localWorkspaces) == 0 {
		return nil
	}

	for i, workspace := range localWorkspaces {
		targetWorkspaceName := TargetNameForWorkspace(monitorDTO.ID, i+1)
		if workspace.Name == targetWorkspaceName {
			continue // noop
		}

		err := HyprctlRenameWorkspace(workspace.ID, targetWorkspaceName)
		if err != nil {
			return fmt.Errorf("rename %q -> %q: %w", workspace.Name, targetWorkspaceName, err)
		}
	}

	return nil
}

func TargetNameForWorkspace(monitorID, workspaceNumber int) string {
	return zeroWidthToken(monitorID) + strconv.Itoa(workspaceNumber)
}

func GetClientsOnWorkspace(workspaceID int, clients []ClientDTO) []ClientDTO {
	result := make([]ClientDTO, 0)
	for _, client := range clients {
		if client.Workspace.ID == workspaceID {
			result = append(result, client)
		}
	}

	return result
}
