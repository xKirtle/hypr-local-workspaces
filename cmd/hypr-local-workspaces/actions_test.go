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
