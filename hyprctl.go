package main

import (
	"fmt"
	"strconv"
)

func runHyprctl(args ...string) error {
	_, status, err := Run("hyprctl", args...)
	if err != nil {
		return err
	}

	if status != 0 {
		return fmt.Errorf("hyprctl %v exited with status %d", args, status)
	}

	return nil
}

func HyprctlWorkspace(name string) error {
	return runHyprctl("dispatch", "workspace", "name:"+name)
}

func HyprctlRenameWorkspace(id int, newName string) error {
	return runHyprctl("dispatch", "renameworkspace", strconv.Itoa(id), newName)
}

func HyprctlMoveToWorkspaceAll(workspaceName string, clients []ClientDTO) error {
	for _, client := range clients {
		err := HyprctlMoveToWorkspace(workspaceName, client.Address)
		if err != nil {
			return err
		}
	}

	return nil
}

func HyprctlMoveToWorkspace(targetName, windowAddr string) error {
	// name:...,address:... must be a single argument
	arg := fmt.Sprintf("name:%s,address:%s", targetName, windowAddr)
	err := runHyprctl("dispatch", "movetoworkspace", arg)
	if err != nil {
		return fmt.Errorf("moving window %q to workspace %q: %w", windowAddr, targetName, err)
	}

	return nil
}

func HyprctlFocusMonitor(monitorId int) error {
	return runHyprctl("dispatch", "focusmonitor", strconv.Itoa(monitorId))
}
