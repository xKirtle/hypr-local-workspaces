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
