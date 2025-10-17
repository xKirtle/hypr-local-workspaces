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

func (m *mockHyprctl) GetClientsInWorkspace(workspaceID int) ([]ClientDTO, error) {
	args := m.Called(workspaceID)
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
