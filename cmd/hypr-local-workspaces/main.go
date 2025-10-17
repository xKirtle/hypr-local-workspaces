package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func parseGotoArgs(args []string) (int, []string, error) {
	if len(args) < 1 {
		return 0, nil, errors.New("usage: hypr-local-workspaces goto <1..9> [global flags]")
	}

	v, err := strconv.Atoi(args[0])
	if err != nil || v < 1 || v > 9 {
		return 0, nil, errors.New("goto index must be a digit 1..9")
	}

	return v, args[1:], nil
}

func parseMoveArgs(args []string) (int, bool, []string, error) {
	fs := flag.NewFlagSet("move", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	all := fs.Bool("all", false, "Apply to all")

	// Separate flags from non-flag args
	if err := fs.Parse(args); err != nil {
		return 0, false, nil, err
	}

	// After parsing, any leftover args are positional
	pos := fs.Args()
	if len(pos) < 1 {
		return 0, false, nil, errors.New("usage: hypr-local-workspaces move <1..9> [--all] [global flags]")
	}

	v, err := strconv.Atoi(pos[0])
	if err != nil {
		return 0, false, nil, errors.New("move expects an integer digit")
	}

	if v < 1 || v > 9 {
		return 0, false, nil, errors.New("move index must be a digit 1..9")
	}

	return v, *all, pos[1:], nil
}

func parseCycleArgs(args []string) (string, []string, error) {
	if len(args) < 1 {
		return "", nil, errors.New("usage: hypr-local-workspaces cycle <next|prev> [global flags]")
	}

	val := strings.ToLower(args[0])
	if val != "next" && val != "prev" {
		return "", nil, errors.New("cycle direction must be 'next' or 'prev'")
	}

	return val, args[1:], nil
}

func parseTrailingGlobalFlags(args []string) (GlobalFlags, error) {
	fs := flag.NewFlagSet("global", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	noCompact := fs.Bool("no-compact", false, "Disable compact mode")

	if err := fs.Parse(args); err != nil {
		return GlobalFlags{Compact: true}, err
	}

	if len(fs.Args()) > 0 {
		return GlobalFlags{Compact: true}, fmt.Errorf("unexpected arguments: %v", fs.Args())
	}

	return GlobalFlags{Compact: !*noCompact}, nil
}

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
