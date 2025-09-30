package main

type MonitorDTO struct {
	ID              int
	Name            string
	Description     string
	Focused         bool
	ActiveWorkspace struct {
		ID   int
		Name string
	}
}

type WorkspaceDTO struct {
	ID        int
	Name      string
	Monitor   string
	MonitorID int
	Windows   int // number of windows/clients
}

type ClientDTO struct {
	Address   string
	Monitor   int
	Workspace struct {
		ID   int
		Name string
	}
}

type Snapshot struct {
	Monitors        []MonitorDTO
	Workspaces      []WorkspaceDTO
	Clients         []ClientDTO
	ActiveWorkspace WorkspaceDTO
	ActiveWindow    ClientDTO
}
