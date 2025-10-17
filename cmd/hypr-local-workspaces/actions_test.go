package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoToWorkspaceDispatchesToRequestedWorkspace(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
		{ID: 4, Name: "1\u200c\u200b", MonitorID: 1},
		{ID: 5, Name: "1\u200d\u200b", MonitorID: 2},
	}, nil)

	dispatcher.On("GoToWorkspace", "3\u200b\u200d").Return(nil)

	action := NewAction(hypr, dispatcher)
	targetIndex := 2
	err := action.GoToWorkspace(targetIndex)

	assert.NoError(t, err)
}

func TestGoToWorkspacePropagatesNoActiveWorkspace(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	hypr.On("GetActiveWorkspace").Return(WorkspaceDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	targetIndex := 2
	err := action.GoToWorkspace(targetIndex)

	assert.Error(t, err)
}

func TestGoToWorkspacePropagatesGetSortedWorkspacesError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	targetIndex := 2
	err := action.GoToWorkspace(targetIndex)

	assert.Error(t, err)
}

func TestGoToWorkspacePropagatesActiveWorkspaceNotOnSortedList(t *testing.T) {
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
		{ID: 4, Name: "1\u200c\u200b", MonitorID: 1},
		{ID: 5, Name: "1\u200d\u200b", MonitorID: 2},
	}, nil)

	action := NewAction(hypr, dispatcher)
	targetIndex := 2
	err := action.GoToWorkspace(targetIndex)

	assert.Error(t, err)
}

func TestGoToWorkspaceReturnsEarlyWhenTargetEqualsCurrent(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
		{ID: 4, Name: "1\u200c\u200b", MonitorID: 1},
		{ID: 5, Name: "1\u200d\u200b", MonitorID: 2},
	}, nil)

	action := NewAction(hypr, dispatcher)
	targetIndex := 1
	err := action.GoToWorkspace(targetIndex)

	assert.NoError(t, err)
}

func TestGoToWorkspacePropagatesGetZeroWidthNameFromIndexError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 999}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 999},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 999},
		{ID: 4, Name: "1\u200c\u200b", MonitorID: 1},
		{ID: 5, Name: "1\u200d\u200b", MonitorID: 2},
	}, nil)

	action := NewAction(hypr, dispatcher)
	targetIndex := 0
	err := action.GoToWorkspace(targetIndex)

	assert.Error(t, err)
}

func TestGoToWorkspacePropagatesDispatcherError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
		{ID: 4, Name: "1\u200c\u200b", MonitorID: 1},
		{ID: 5, Name: "1\u200d\u200b", MonitorID: 2},
	}, nil)

	dispatcher.On("GoToWorkspace", "3\u200b\u200d").Return(assert.AnError)

	action := NewAction(hypr, dispatcher)
	targetIndex := 2
	err := action.GoToWorkspace(targetIndex)

	assert.Error(t, err)
}

// func TestDummyTest(t *testing.T) {
// 	hypr := NewHyprctlClient(2)
// 	dispatcher := NewDispatcherClient()
// 	action := NewAction(hypr, dispatcher)

// 	_ = action.GoToWorkspace(1)
// }
