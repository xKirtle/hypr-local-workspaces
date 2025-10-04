package main

import "time"

type MonitorDTO struct {
	ID              int
	Name            string
	Focused         bool
	ActiveWorkspace SimpleWorkspace
}

type WorkspaceDTO struct {
	ID           int
	Name         string
	Monitor      string
	MonitorID    int
	WindowsCount int `json:"windows"`
}

type ClientDTO struct {
	Address   string
	Monitor   int
	Workspace SimpleWorkspace
}

type SimpleWorkspace struct {
	ID   int
	Name string
}

type Action struct {
	hyprctl    hyprctl
	dispatcher dispatcher
}

type hyprctl interface {
	GetMonitors() ([]MonitorDTO, error)
	GetWorkspaces() ([]WorkspaceDTO, error)
	GetClients() ([]ClientDTO, error)
	GetActiveWorkspace() (WorkspaceDTO, error)
	GetActiveWindow() (ClientDTO, error)
	GetActiveMonitorID() (int, error)
}

type dispatcher interface {
	Workspace(wsName string) error
	RenameWorkspace(id int, wsNewName string) error
	FocusMonitor(monitorId int) error
	MoveAllToWorkspace(wsName string) error
	MoveToWorkspace(wsName, windowAddr string) error
}

type hyprctlClient struct {
	timeout time.Duration
}

type dispatcherClient struct{}
