package main

import (
	"fmt"
	"strconv"
)

func NewDispatcherClient() dispatcher {
	return &dispatcherClient{}
}

func hyprDispatch(args ...string) error {
	allArgs := append([]string{"dispatch"}, args...)
	_, _, err := RunWith("hyprctl", allArgs, CaptureOutput(), WithTimeout(HyprctlTimeout))

	if err != nil {
		return err
	}

	return nil
}

func (d *dispatcherClient) GoToWorkspace(wsName string) error {
	return hyprDispatch("workspace", fmt.Sprintf("name:%s", wsName))
}

func (d *dispatcherClient) RenameWorkspace(id int, wsNewName string) error {
	return hyprDispatch("renameworkspace", strconv.Itoa(id), wsNewName)
}

func (d *dispatcherClient) FocusMonitor(monitorId int) error {
	return hyprDispatch("focusmonitor", strconv.Itoa(monitorId))
}

func (d *dispatcherClient) MoveToWorkspace(wsName string) error {
	return hyprDispatch("movetoworkspace", wsName)
}

func (d *dispatcherClient) MoveAddrToWorkspace(wsName, windowAddr string) error {
	return hyprDispatch("movetoworkspace", fmt.Sprintf("name:%s,address:%s", wsName, windowAddr))
}
