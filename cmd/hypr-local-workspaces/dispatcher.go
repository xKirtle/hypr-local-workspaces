package main

import (
	"fmt"
	"strconv"
	"strings"
)

func NewDispatcherClient() dispatcher {
	return &dispatcherClient{}
}

func hyprDispatchLua(expr string) error {
    output, exitCode, err := RunWith(
        "hyprctl",
        []string{"dispatch", expr},
        CaptureOutput(),
        CombineOutput(),
        WithTimeout(HyprctlTimeout),
    )

    if err != nil {
        return err
    }

    if exitCode != 0 {
        return fmt.Errorf(
            "hyprctl dispatch failed with exit code %d: %s",
            exitCode,
            strings.TrimSpace(string(output)),
        )
    }

    return nil
}

func (d *dispatcherClient) GoToWorkspace(wsName string) error {
    expr := fmt.Sprintf(
        `hl.dsp.focus({
            workspace = %s,
            on_current_monitor = true
        })`,
        luaString("name:"+wsName),
    )

    return hyprDispatchLua(expr)
}

func (d *dispatcherClient) RenameWorkspace(id int, wsNewName string) error {
    expr := fmt.Sprintf(
        `hl.dsp.workspace.rename({ workspace = %d, name = %s })`,
        id,
        luaString(wsNewName),
    )

    return hyprDispatchLua(expr)
}

func (d *dispatcherClient) FocusMonitor(monitorId int) error {
    expr := fmt.Sprintf(
        `hl.dsp.focus({ monitor = %s })`,
        luaString(strconv.Itoa(monitorId)),
    )

    return hyprDispatchLua(expr)
}

func (d *dispatcherClient) MoveToWorkspace(wsName string) error {
    expr := fmt.Sprintf(
        `hl.dsp.window.move({
            workspace = %s,
            follow = true
        })`,
        luaString("name:" + wsName),
    )

    return hyprDispatchLua(expr)
}

func (d *dispatcherClient) MoveAddrToWorkspace(
    wsName,
    windowAddr string,
) error {
    expr := fmt.Sprintf(
        `hl.dsp.window.move({
            workspace = %s,
            window = %s,
            follow = false
        })`,
        luaString("name:" + wsName),
        luaString("address:" + windowAddr),
    )

    return hyprDispatchLua(expr)
}

func luaString(s string) string {
    // Lua long-bracket strings preserve UTF-8 bytes verbatim and require
    // no escaping. Increase the '=' level if the closing delimiter happens
    // to occur in the input.
    for level := 0; ; level++ {
        equals := strings.Repeat("=", level)
        close := "]" + equals + "]"

        if !strings.Contains(s, close) {
            return "[" + equals + "[" + s + "]" + equals + "]"
        }
    }
}
