package main

import "github.com/stretchr/testify/mock"

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
