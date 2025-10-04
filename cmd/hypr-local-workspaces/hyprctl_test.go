package main

import "github.com/stretchr/testify/mock"

type mockHyprctl struct {
	mock.Mock
}

func (m *mockHyprctl) GetMonitors() ([]MonitorDTO, error) {
	args := m.Called()
	ms, _ := args.Get(0).([]MonitorDTO)
	return ms, args.Error(1)
}

func (m *mockHyprctl) GetWorkspaces() ([]WorkspaceDTO, error) {
	args := m.Called()
	ws, _ := args.Get(0).([]WorkspaceDTO)
	return ws, args.Error(1)
}

func (m *mockHyprctl) GetClients() ([]ClientDTO, error) {
	args := m.Called()
	cs, _ := args.Get(0).([]ClientDTO)
	return cs, args.Error(1)
}

func (m *mockHyprctl) GetActiveWorkspace() (WorkspaceDTO, error) {
	args := m.Called()
	ws, _ := args.Get(0).(WorkspaceDTO)
	return ws, args.Error(1)
}

func (m *mockHyprctl) GetActiveWindow() (ClientDTO, error) {
	args := m.Called()
	c, _ := args.Get(0).(ClientDTO)
	return c, args.Error(1)
}

func (m *mockHyprctl) GetActiveMonitorID() (int, error) {
	args := m.Called()
	id, _ := args.Get(0).(int)
	return id, args.Error(1)
}

type mockDispatcher struct {
	mock.Mock
}

func (m *mockDispatcher) Workspace(wsName string) error {
	args := m.Called(wsName)
	return args.Error(0)
}

func (m *mockDispatcher) RenameWorkspace(id int, wsNewName string) error {
	args := m.Called(id, wsNewName)
	return args.Error(0)
}

func (m *mockDispatcher) FocusMonitor(monitorId int) error {
	args := m.Called(monitorId)
	return args.Error(0)
}

func (m *mockDispatcher) MoveAllToWorkspace(wsName string) error {
	args := m.Called(wsName)
	return args.Error(0)
}

func (m *mockDispatcher) MoveToWorkspace(wsName, windowAddr string) error {
	args := m.Called(wsName, windowAddr)
	return args.Error(0)
}

// TODO: Maybe I don't need to recreate a live hyprctl state for tests, just mock the calls and return values I expect.
func sampleMonitors() []MonitorDTO {
	return []MonitorDTO{
		{
			ID:      0,
			Name:    "DP-1",
			Focused: true,
			ActiveWorkspace: SimpleWorkspace{
				ID:   1,
				Name: "1\u200b\u200b",
			},
		},
		{
			ID:      1,
			Name:    "DP-2",
			Focused: false,
			ActiveWorkspace: SimpleWorkspace{
				ID:   2,
				Name: "1\u200c\u200b",
			},
		},
		{
			ID:      2,
			Name:    "DP-3",
			Focused: false,
			ActiveWorkspace: SimpleWorkspace{
				ID:   3,
				Name: "1\u200d\u200b",
			},
		},
	}
}
