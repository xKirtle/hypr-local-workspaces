package main

import "time"

func NewHyprctlClient(timeout time.Duration) hyprctl {
	return &hyprctlClient{timeout: timeout}
}

func NewDispatcherClient() dispatcher {
	return &dispatcherClient{}
}

func (c *hyprctlClient) GetMonitors() ([]MonitorDTO, error) {
	return nil, nil
}

func (c *hyprctlClient) GetWorkspaces() ([]WorkspaceDTO, error) {
	return nil, nil
}

func (c *hyprctlClient) GetClients() ([]ClientDTO, error) {
	return nil, nil
}

func (c *hyprctlClient) GetActiveWorkspace() (WorkspaceDTO, error) {
	return WorkspaceDTO{}, nil
}

func (c *hyprctlClient) GetActiveWindow() (ClientDTO, error) {
	return ClientDTO{}, nil
}

func (c *hyprctlClient) GetActiveMonitorID() (int, error) {
	return 0, nil
}

func (d *dispatcherClient) Workspace(wsName string) error {
	return nil
}

func (d *dispatcherClient) RenameWorkspace(id int, wsNewName string) error {
	return nil
}

func (d *dispatcherClient) FocusMonitor(monitorId int) error {
	return nil
}

func (d *dispatcherClient) MoveAllToWorkspace(wsName string) error {
	return nil
}

func (d *dispatcherClient) MoveToWorkspace(wsName, windowAddr string) error {
	return nil
}
