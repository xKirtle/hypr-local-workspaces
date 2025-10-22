package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}

	subcmd := args[0]
	subArgs := args[1:]

	action := NewAction(
		NewHyprctlClient(2*time.Second),
		NewDispatcherClient(),
	)

	switch subcmd {
	case "goto":
		targetWorkspace, trailing, err := parseGotoArgs(subArgs)
		if err != nil {
			fail(err)
		}

		globals, err := parseTrailingGlobalFlags(trailing)
		if err != nil {
			fail(err)
		}

		targetIndex := targetWorkspace - 1
		_ = action.GoToWorkspace(targetIndex, globals.Compact)

	case "move":
		targetWorkspace, all, trailing, err := parseMoveArgs(subArgs)
		if err != nil {
			fail(err)
		}

		globals, err := parseTrailingGlobalFlags(trailing)
		if err != nil {
			fail(err)
		}

		targetIndex := targetWorkspace - 1
		_ = action.MoveToWorkspace(targetIndex, all, globals.Compact)

	case "cycle":
		dir, trailing, err := parseCycleArgs(subArgs)
		if err != nil {
			fail(err)
		}

		globals, err := parseTrailingGlobalFlags(trailing)
		if err != nil {
			fail(err)
		}

		_ = action.CycleWorkspace(dir, globals.Compact)

	case "init":
		_ = action.InitWorkspaces()

	case "help", "-h", "--help", "":
		printUsage()

	default:
		fail(fmt.Errorf("unknown subcommand: %q", subcmd))
	}
}

// parsing helpers moved to parse.go

func printUsage() {
	_, _ = fmt.Fprintln(os.Stderr, `Usage:
  hypr-local-workspaces goto  <1..9>         [global flags]
  hypr-local-workspaces move  <1..9> [--all] [global flags]
  hypr-local-workspaces cycle <next|prev>    [global flags]

Global flags:
  --no-compact    Disable compact mode (enabled by default)`)
}

func fail(err error) {
	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	printUsage()
	os.Exit(ExitMissingArgs)
}

// Could probably rely on a third-party library that supports nested subcommands, persistent global flags, etc...
// but this is sufficient for now.
