package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCycleWorkspace_NoActiveWorkspaceError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	hypr.On("GetActiveWorkspace").Return(WorkspaceDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", true)
	assert.Error(t, err)
}

func TestCycleWorkspace_GetWorkspacesError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", true)
	assert.Error(t, err)
}

func TestCycleWorkspace_CurrentNotOnListError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 42, Name: "42\u200b\u200c", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", true)
	assert.Error(t, err)
}

func TestCycleWorkspace_ReturnsEarlyWhenPrevOnFirst(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
	}, nil)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("prev", true)
	assert.NoError(t, err)
}

func TestCycleWorkspace_CompactTrue_Success(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)

	// Cycling next from index 0 -> index 1
	dispatcher.On("GoToWorkspace", "2\u200b\u200c").Return(nil)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", true)
	assert.NoError(t, err)
}

func TestCycleWorkspace_DispatcherError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
	}, nil)

	dispatcher.On("GoToWorkspace", "2\u200b\u200c").Return(assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", true)
	assert.Error(t, err)
}

func TestCycleWorkspace_CompactFalse_UsesExistingName(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
	}, nil)

	dispatcher.On("GoToWorkspace", "2\u200b\u200c").Return(nil)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", false)
	assert.NoError(t, err)
}

func TestCycleWorkspace_ZeroWidthNameError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 999}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 999},
	}, nil)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", true)
	assert.Error(t, err)
}

func TestCycleWorkspace_CompactionError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	// First call for sorted list
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
	}, nil).Once()
	// Second call for compaction fails
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", true)
	assert.Error(t, err)
}

func TestCycleWorkspace_CompactFalse_DispatcherError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
	}, nil)

	dispatcher.On("GoToWorkspace", "2\u200b\u200c").Return(assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.CycleWorkspace("next", false)
	assert.Error(t, err)
}
