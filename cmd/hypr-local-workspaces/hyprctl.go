package main

import (
	"encoding/json"
	"time"
)

const (
	HyprctlTimeout = 2 * time.Second
)

func NewHyprctlClient(timeout time.Duration) hyprctl {
	return &hyprctlClient{timeout: timeout}
}

func hyprJson(cmd string) ([]byte, error) {
	args := []string{"-j", cmd}
	out, _, err := RunWith("hyprctl", args, CaptureOutput(), WithTimeout(HyprctlTimeout))

	if err != nil {
		return nil, err
	}

	return out, nil
}

func hyprJsonDecode[T any](cmd string) (T, error) {
	out, err := hyprJson(cmd)

	if err != nil {
		var emptyT T
		return emptyT, err
	}

	var result T
	err = json.Unmarshal(out, &result)

	if err != nil {
		var emptyT T
		return emptyT, err
	}

	return result, nil
}

func (c *hyprctlClient) GetMonitors() ([]MonitorDTO, error) {
	return hyprJsonDecode[[]MonitorDTO]("monitors")
}

func (c *hyprctlClient) GetWorkspaces() ([]WorkspaceDTO, error) {
	return hyprJsonDecode[[]WorkspaceDTO]("workspaces")
}

func (c *hyprctlClient) GetClients() ([]ClientDTO, error) {
	return hyprJsonDecode[[]ClientDTO]("clients")
}

func (c *hyprctlClient) GetClientsInWorkspace(workspaceID int) ([]ClientDTO, error) {
	clients, err := c.GetClients()
	if err != nil {
		return nil, err
	}

	var filtered []ClientDTO
	for _, client := range clients {
		if client.Workspace.ID == workspaceID {
			filtered = append(filtered, client)
		}
	}

	return filtered, nil
}

func (c *hyprctlClient) GetActiveWorkspace() (WorkspaceDTO, error) {
	return hyprJsonDecode[WorkspaceDTO]("activeworkspace")
}

func (c *hyprctlClient) GetActiveWindow() (ClientDTO, error) {
	return hyprJsonDecode[ClientDTO]("activewindow")
}

func (c *hyprctlClient) GetActiveMonitorID() (int, error) {
	activeWs, err := c.GetActiveWorkspace()

	if err != nil {
		return -1, err
	}

	return activeWs.MonitorID, nil
}
