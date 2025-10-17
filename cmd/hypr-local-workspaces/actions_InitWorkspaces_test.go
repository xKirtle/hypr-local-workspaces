package main

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestInitWorkspaces_GetMonitorsError(t *testing.T) {
    hypr := new(mockHyprctl)
    dispatcher := new(mockDispatcher)
    defer hypr.AssertExpectations(t)
    defer dispatcher.AssertExpectations(t)

    hypr.On("GetMonitors").Return([]MonitorDTO{}, assert.AnError)

    action := NewAction(hypr, dispatcher)
    err := action.InitWorkspaces()
    assert.Error(t, err)
}

func TestInitWorkspaces_Success(t *testing.T) {
    hypr := new(mockHyprctl)
    dispatcher := new(mockDispatcher)
    defer hypr.AssertExpectations(t)
    defer dispatcher.AssertExpectations(t)

    hypr.On("GetMonitors").Return([]MonitorDTO{{ID: 0}, {ID: 1}}, nil)
    // Workspaces for both monitors, correctly named and will not trigger renames
    hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
        {ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
        {ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
        {ID: 3, Name: "1\u200c\u200b", MonitorID: 1},
        {ID: 4, Name: "2\u200c\u200c", MonitorID: 1},
    }, nil)

    action := NewAction(hypr, dispatcher)
    err := action.InitWorkspaces()
    assert.NoError(t, err)
}

func TestInitWorkspaces_CompactionError(t *testing.T) {
    hypr := new(mockHyprctl)
    dispatcher := new(mockDispatcher)
    defer hypr.AssertExpectations(t)
    defer dispatcher.AssertExpectations(t)

    hypr.On("GetMonitors").Return([]MonitorDTO{{ID: 0}}, nil)
    // Cause compaction to fail by returning error on GetWorkspaces
    hypr.On("GetWorkspaces").Return([]WorkspaceDTO{}, assert.AnError)

    action := NewAction(hypr, dispatcher)
    err := action.InitWorkspaces()
    assert.Error(t, err)
}

