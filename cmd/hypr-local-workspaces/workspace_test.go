package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWorkspacesOnMonitorFiltersByMonitorID(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	monitorID := 1
	expected := []WorkspaceDTO{
		{ID: 1, Name: "ws-1", MonitorID: monitorID},
		{ID: 3, Name: "ws-3", MonitorID: monitorID},
	}

	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		expected[0],
		{ID: 2, Name: "ws-2", MonitorID: 7},
		expected[1],
	}, nil)

	workspaces, err := GetWorkspacesOnMonitor(hypr, monitorID)

	require.NoError(t, err)
	assert.Equal(t, expected, workspaces)
}

func TestGetWorkspacesOnMonitorReturnsEmptyList(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	monitorID := 1
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 2, Name: "ws-2", MonitorID: 7},
	}, nil)

	workspaces, err := GetWorkspacesOnMonitor(hypr, monitorID)

	require.NoError(t, err)
	assert.Empty(t, workspaces)
}

func TestGetWorkspacesOnMonitorPropagatesErrors(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	sentinelErr := errors.New("get workspaces failed")
	hypr.On("GetWorkspaces").Return(([]WorkspaceDTO)(nil), sentinelErr)

	workspaces, err := GetWorkspacesOnMonitor(hypr, 42)

	assert.Nil(t, workspaces)
	assert.ErrorIs(t, err, sentinelErr)
}

func TestGetSortedWorkspacesOnMonitorSortsByNameIgnoringZeroWidthChars(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	monitorID := 1
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 3, Name: "3\u200b\u200d", MonitorID: monitorID},
		{ID: 1, Name: "1\u200c\u200b", MonitorID: monitorID},
		{ID: 5, Name: "5\u200f\u2064", MonitorID: monitorID},
		{ID: 4, Name: "10\u200b\u200c\u200c", MonitorID: monitorID},
		{ID: 2, Name: "6\u200f\u2060", MonitorID: monitorID},
	}, nil)

	expected := []WorkspaceDTO{
		{ID: 1, Name: "1\u200c\u200b", MonitorID: monitorID},
		{ID: 3, Name: "3\u200b\u200d", MonitorID: monitorID},
		{ID: 5, Name: "5\u200f\u2064", MonitorID: monitorID},
		{ID: 2, Name: "6\u200f\u2060", MonitorID: monitorID},
		{ID: 4, Name: "10\u200b\u200c\u200c", MonitorID: monitorID},
	}

	workspaces, err := GetSortedWorkspacesOnMonitor(hypr, monitorID)

	require.NoError(t, err)
	assert.Equal(t, expected, workspaces)
}

func TestGetSortedWorkspacesOnMonitorPropagatesErrorsFromGetWorkspacesOnMonitor(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	sentinelErr := errors.New("get workspaces failed")
	hypr.On("GetWorkspaces").Return(([]WorkspaceDTO)(nil), sentinelErr)

	workspaces, err := GetSortedWorkspacesOnMonitor(hypr, 42)

	assert.Nil(t, workspaces)
	assert.ErrorIs(t, err, sentinelErr)
}

func TestGetSortedWorkspacesOnMonitorSortsByIdWhenInvalidZeroWidthNames(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	monitorID := 1
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 3, Name: "ws-1", MonitorID: monitorID},
		{ID: 1, Name: "ws-2", MonitorID: monitorID},
		{ID: 6, Name: "invalid-\xff", MonitorID: monitorID}, // Invalid UTF-8
		{ID: 5, Name: "ws-3", MonitorID: monitorID},
		{ID: 4, Name: "ws-\u200b3", MonitorID: monitorID},
		{ID: 2, Name: "ws-10", MonitorID: monitorID},
	}, nil)

	expected := []WorkspaceDTO{
		{ID: 1, Name: "ws-2", MonitorID: monitorID},
		{ID: 2, Name: "ws-10", MonitorID: monitorID},
		{ID: 3, Name: "ws-1", MonitorID: monitorID},
		{ID: 4, Name: "ws-\u200b3", MonitorID: monitorID},
		{ID: 5, Name: "ws-3", MonitorID: monitorID},
		{ID: 6, Name: "invalid-\xff", MonitorID: monitorID}, // Should be last due to ID fallback
	}

	workspaces, err := GetSortedWorkspacesOnMonitor(hypr, monitorID)

	require.NoError(t, err)
	assert.Equal(t, expected, workspaces)
}

func TestGetSortedWorkspacesOnMonitorReturnsEmptyList(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	monitorID := 1
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 2, Name: "ws-2", MonitorID: 7},
	}, nil)

	workspaces, err := GetSortedWorkspacesOnMonitor(hypr, monitorID)

	require.NoError(t, err)
	assert.Empty(t, workspaces)
}

func TestGetSortedWorkspacesOnMonitorSortsByIdWhenInvalidNames(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	monitorID := 1
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 3, Name: "", MonitorID: monitorID}, // Invalid UTF-8
		{ID: 1, Name: "", MonitorID: monitorID}, // Invalid UTF-8
	}, nil)

	expected := []WorkspaceDTO{
		{ID: 1, Name: "", MonitorID: monitorID},
		{ID: 3, Name: "", MonitorID: monitorID},
	}

	workspaces, err := GetSortedWorkspacesOnMonitor(hypr, monitorID)

	require.NoError(t, err)
	assert.Equal(t, expected, workspaces)
}

func TestGetSortedWorkspacesOnMonitorGuardsAgainstAtoiOverflow(t *testing.T) {
	hypr := new(mockHyprctl)
	defer hypr.AssertExpectations(t)

	monitorID := 1
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 3, Name: "99999999999999999999999999999999999999999999999999", MonitorID: monitorID}, // Overflow
		{ID: 1, Name: "1", MonitorID: monitorID},
	}, nil)

	expected := []WorkspaceDTO{
		{ID: 1, Name: "1", MonitorID: monitorID},
		{ID: 3, Name: "99999999999999999999999999999999999999999999999999", MonitorID: monitorID},
	}

	workspaces, err := GetSortedWorkspacesOnMonitor(hypr, monitorID)

	require.NoError(t, err)
	assert.Equal(t, expected, workspaces)
}

func TestGetWorkspaceIndexOnList_FindsIndex(t *testing.T) {
	list := []WorkspaceDTO{
		{ID: 10, Name: "10"},
		{ID: 2, Name: "2"},
		{ID: 5, Name: "5"},
	}

	idx := GetWorkspaceIndexOnList(list, 2)
	assert.Equal(t, 1, idx)

	idx = GetWorkspaceIndexOnList(list, 10)
	assert.Equal(t, 0, idx)

	idx = GetWorkspaceIndexOnList(list, 5)
	assert.Equal(t, 2, idx)
}

func TestGetWorkspaceIndexOnList_NotFound(t *testing.T) {
	list := []WorkspaceDTO{
		{ID: 1, Name: "1"},
		{ID: 3, Name: "3"},
	}

	idx := GetWorkspaceIndexOnList(list, 2)
	assert.Equal(t, -1, idx)
}

func TestGetWorkspaceIndexOnList_EmptyList(t *testing.T) {
	var list []WorkspaceDTO
	idx := GetWorkspaceIndexOnList(list, 1)
	assert.Equal(t, -1, idx)
}

func TestCompactLocalWorkspacesOnMonitor_CompactsList(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	monitorID := 0
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
		{ID: 10, Name: "6\u200b\u2061", MonitorID: 0},
		{ID: 15, Name: "8\u200b\u2063", MonitorID: 0},
		{ID: 4, Name: "1\u200c\u200b", MonitorID: 1},
		{ID: 5, Name: "1\u200d\u200b", MonitorID: 2},
	}, nil)

	expected := []WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 3, Name: "2\u200b\u200c", MonitorID: 0},
		{ID: 10, Name: "3\u200b\u200d", MonitorID: 0},
		{ID: 15, Name: "4\u200b\u200e", MonitorID: 0},
	}

	dispatcher.On("RenameWorkspace", 3, expected[1].Name).Return(nil)
	dispatcher.On("RenameWorkspace", 10, expected[2].Name).Return(nil)
	dispatcher.On("RenameWorkspace", 15, expected[3].Name).Return(nil)

	action := &Action{hyprctl: hypr, dispatcher: dispatcher}
	err := CompactLocalWorkspacesOnMonitor(action, monitorID)

	assert.NoError(t, err)
}

func TestCompactLocalWorkspacesOnMonitor_PropagatesErrors(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	monitorID := 0
	sentinelErr := errors.New("get workspaces failed")
	hypr.On("GetWorkspaces").Return(([]WorkspaceDTO)(nil), sentinelErr)

	action := &Action{hyprctl: hypr, dispatcher: dispatcher}
	err := CompactLocalWorkspacesOnMonitor(action, monitorID)

	assert.ErrorIs(t, err, sentinelErr)
}

func TestCompactLocalWorkspacesOnMonitor_PropagatesGetZeroWidthNameToIndexErrors(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	monitorID := 0
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 3, Name: "invalid-\xff", MonitorID: 0}, // Invalid UTF-8
	}, nil)

	action := &Action{hyprctl: hypr, dispatcher: dispatcher}
	err := CompactLocalWorkspacesOnMonitor(action, monitorID)

	assert.Error(t, err)
}

func TestCompactLocalWorkspacesOnMonitor_HandlesNoRenameNeeded(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	monitorID := 0
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 3, Name: "2\u200b\u200c", MonitorID: 0},
		{ID: 10, Name: "3\u200b\u200d", MonitorID: 0},
		{ID: 15, Name: "4\u200b\u200e", MonitorID: 0},
	}, nil)

	action := &Action{hyprctl: hypr, dispatcher: dispatcher}
	err := CompactLocalWorkspacesOnMonitor(action, monitorID)

	assert.NoError(t, err)
}

func TestCompactLocalWorkspacesOnMonitor_PropagatesRenameWorkspaceErrors(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	monitorID := 0
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
		{ID: 10, Name: "6\u200b\u2061", MonitorID: 0},
	}, nil)

	dispatcher.On("RenameWorkspace", 3, "2\u200b\u200c").Return(nil)
	sentinelErr := errors.New("rename failed")
	dispatcher.On("RenameWorkspace", 10, "3\u200b\u200d").Return(sentinelErr)

	action := &Action{hyprctl: hypr, dispatcher: dispatcher}
	err := CompactLocalWorkspacesOnMonitor(action, monitorID)

	assert.ErrorIs(t, err, sentinelErr)
}
