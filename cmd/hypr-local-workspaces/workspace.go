package main

import "sort"

func GetWorkspacesOnMonitor(hyprctl hyprctl, monitorId int) ([]WorkspaceDTO, error) {
	workspaces, err := hyprctl.GetWorkspaces()
	if err != nil {
		return nil, err
	}

	var monitorWorkspaces []WorkspaceDTO
	for _, ws := range workspaces {
		if ws.MonitorID == monitorId {
			monitorWorkspaces = append(monitorWorkspaces, ws)
		}
	}

	return monitorWorkspaces, nil
}

func GetSortedWorkspacesOnMonitor(hyprctl hyprctl, monitorId int) ([]WorkspaceDTO, error) {
	workspaces, err := GetWorkspacesOnMonitor(hyprctl, monitorId)
	if err != nil {
		return nil, err
	}

	// Sort by name, ignoring zero-width chars
	sort.Slice(workspaces, func(i, j int) bool {
		nameI, errI := GetZeroWidthNameToIndex(workspaces[i].Name)
		nameJ, errJ := GetZeroWidthNameToIndex(workspaces[j].Name)
		if errI != nil || errJ != nil {
			// TODO: How to handle errors here?
			// For now, just fall back to id comparison

			return workspaces[i].ID < workspaces[j].ID
		}

		return nameI < nameJ
	})

	return workspaces, nil
}

// DecideTargetWorkspaceIndex normalizes a requested index into [0..N], where N == len(sortedWorkspaces)
// represents a request to create a new workspace slot. It also returns a compact flag indicating
// whether moving would leave gaps (e.g., leaving or skipping over an empty workspace).
// This function is intentionally unopinionated about no-ops. Callers should compare the returned
// index to currentIndex and decide if it is a no-op.
func DecideTargetWorkspaceIndex(currentIndex, requested int, sortedWorkspaces []WorkspaceDTO) (int, bool) {
	n := len(sortedWorkspaces)

	// Normalize into [0..n]
	targetIndex := requested
	if targetIndex < 0 {
		targetIndex = 0
	} else if targetIndex > n {
		targetIndex = n
	}

	// Special-case: avoid compaction when on last empty and requesting a new slot.
	if n > 0 && targetIndex == n && currentIndex == n-1 {
		if sortedWorkspaces[currentIndex].WindowsCount == 0 {
			return currentIndex, false
		}
	}

	compact := false

	// Only compute compaction when currentIndex is within existing bounds and we're moving.
	if targetIndex != currentIndex && currentIndex >= 0 && currentIndex < n {
		// Leaving an empty workspace will require compaction.
		if sortedWorkspaces[currentIndex].WindowsCount == 0 {
			compact = true
		}

		// Check if we're skipping over any empty existing workspace between current and target.
		low, high := currentIndex, targetIndex
		if low > high {
			low, high = targetIndex, currentIndex
		}
		if high > n-1 {
			high = n - 1 // exclude the synthetic N index from scanning
		}
		for i := low + 1; i <= high; i++ {
			if i >= 0 && i < n && sortedWorkspaces[i].WindowsCount == 0 {
				compact = true
				break
			}
		}
	}

	return targetIndex, compact
}

func GetWorkspaceIndexOnList(sortedLocalWs []WorkspaceDTO, workspaceID int) int {
	for i, ws := range sortedLocalWs {
		if ws.ID == workspaceID {
			return i
		}
	}

	return -1
}

func CompactLocalWorkspacesOnMonitor(action *Action, monitorID int) error {
	hyprctl, dispatcher := action.hyprctl, action.dispatcher

	sortedLocalWs, err := GetSortedWorkspacesOnMonitor(hyprctl, monitorID)
	if err != nil {
		return err
	}

	for i, ws := range sortedLocalWs {
		wsIndex, err := GetZeroWidthNameToIndex(ws.Name)

		if err != nil {
			return err
		}

		if wsIndex == i {
			continue
		}

		newName, err := GetZeroWidthNameFromIndex(monitorID, i)

		// Can't really happen either? monitorID or index i would have to be out of range
		// However, monitorID is also checked when fetching sortedLocalWs above
		// So really only index i would have to be out of range, which is impossible in this loop?
		if err != nil {
			return err
		}

		// Should never happen
		if ws.Name == newName {
			continue
		}

		if err := dispatcher.RenameWorkspace(ws.ID, newName); err != nil {
			return err
		}
	}

	return nil
}
