package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveToWorkspace_NoActiveWorkspaceError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	hypr.On("GetActiveWorkspace").Return(WorkspaceDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(1, false, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_GetWorkspacesError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, false, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_CurrentNotOnListError(t *testing.T) {
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
	err := action.MoveToWorkspace(1, false, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_ReturnsEarlyWhenTargetEqualsCurrent(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0, WindowsCount: 1}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)

	action := NewAction(hypr, dispatcher)
	// targetIndex points to current workspace (index 1)
	err := action.MoveToWorkspace(1, false, true)
	assert.NoError(t, err)
}

func TestMoveToWorkspace_ReturnsEarlyWhenSingleWindowMovingForward(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 3, Name: "3\u200b\u200d", MonitorID: 0, WindowsCount: 1}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
		activeWs,
	}, nil)

	action := NewAction(hypr, dispatcher)
	// targetIndex points beyond current workspace (index 3)
	err := action.MoveToWorkspace(3, false, true)
	assert.NoError(t, err)
}

func TestMoveToWorkspace_ZeroWidthNameError(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 999, WindowsCount: 3}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 999},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 999},
	}, nil)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, false, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_MoveAllClients_Success_WithCompaction(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0, WindowsCount: 3}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)
	hypr.On("GetClientsInWorkspace", activeWs.ID).Return([]ClientDTO{
		{Address: "0xabc"}, {Address: "0xdef"},
	}, nil)

	// Target index 2 -> target name for monitor 0 index 2
	dispatcher.On("MoveAddrToWorkspace", "3\u200b\u200d", "0xabc").Return(nil)
	dispatcher.On("MoveAddrToWorkspace", "3\u200b\u200d", "0xdef").Return(nil)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, true, true)
	assert.NoError(t, err)
}

func TestMoveToWorkspace_MoveAllClients_Error_GetClients(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0, WindowsCount: 3}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)
	hypr.On("GetClientsInWorkspace", activeWs.ID).Return([]ClientDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, true, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_MoveAllClients_Error_MoveAddr(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0, WindowsCount: 3}

	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)
	hypr.On("GetClientsInWorkspace", activeWs.ID).Return([]ClientDTO{
		{Address: "0xabc"}, {Address: "0xdef"},
	}, nil)

	dispatcher.On("MoveAddrToWorkspace", "3\u200b\u200d", "0xabc").Return(assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, true, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_MoveSingleClient_Success_WithCompaction(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0, WindowsCount: 3}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)
	hypr.On("GetActiveWindow").Return(ClientDTO{Address: "0xabc"}, nil)

	dispatcher.On("MoveAddrToWorkspace", "3\u200b\u200d", "0xabc").Return(nil)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, false, true)
	assert.NoError(t, err)
}

func TestMoveToWorkspace_MoveSingleClient_Error_GetActiveWindow(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0, WindowsCount: 3}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)
	hypr.On("GetActiveWindow").Return(ClientDTO{}, assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, false, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_MoveSingleClient_Error_MoveTo(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 2, Name: "2\u200b\u200c", MonitorID: 0, WindowsCount: 3}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		activeWs,
		{ID: 3, Name: "3\u200b\u200d", MonitorID: 0},
	}, nil)
	hypr.On("GetActiveWindow").Return(ClientDTO{Address: "0xabc"}, nil)

	dispatcher.On("MoveAddrToWorkspace", "3\u200b\u200d", "0xabc").Return(assert.AnError)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(2, false, true)
	assert.Error(t, err)
}

func TestMoveToWorkspace_NoCompactionWhenCompactFalse(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 1, Name: "1\u200b\u200b", MonitorID: 0, WindowsCount: 3}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		activeWs,
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
	}, nil)
	hypr.On("GetActiveWindow").Return(ClientDTO{Address: "0xabc"}, nil)

	dispatcher.On("MoveAddrToWorkspace", "2\u200b\u200c", "0xabc").Return(nil)

	action := NewAction(hypr, dispatcher)
	// compact false should skip compaction path
	err := action.MoveToWorkspace(1, false, false)
	assert.NoError(t, err)
}

func TestMoveToWorkspace_CompactionErrorAfterMove(t *testing.T) {
	hypr := new(mockHyprctl)
	dispatcher := new(mockDispatcher)
	defer hypr.AssertExpectations(t)
	defer dispatcher.AssertExpectations(t)

	activeWs := WorkspaceDTO{ID: 3, Name: "3\u200b\u200d", MonitorID: 0, WindowsCount: 1}
	hypr.On("GetActiveWorkspace").Return(activeWs, nil)
	// First call (for sorting) succeeds
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{
		{ID: 1, Name: "1\u200b\u200b", MonitorID: 0},
		{ID: 2, Name: "2\u200b\u200c", MonitorID: 0},
		activeWs,
	}, nil).Once()
	// Second call (compaction) fails
	hypr.On("GetWorkspaces").Return([]WorkspaceDTO{}, assert.AnError)
	hypr.On("GetActiveWindow").Return(ClientDTO{Address: "0xabc"}, nil)

	dispatcher.On("MoveAddrToWorkspace", "1\u200b\u200b", "0xabc").Return(nil)

	action := NewAction(hypr, dispatcher)
	err := action.MoveToWorkspace(0, false, true)
	assert.Error(t, err)
}
