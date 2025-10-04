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
	subcmd, subArgs, err := splitSubcommand(os.Args)
	if err != nil {
		fail(err)
	}

	action := NewAction(
		NewHyprctlClient(2*time.Second),
		NewDispatcherClient(),
	)

	switch subcmd {
	case "goto":
		targetWorkspace, err := parseGotoArgs(subArgs)
		if err != nil {
			fail(err)
		}

		targetIndex := targetWorkspace - 1
		_ = action.GoToWorkspace(targetIndex)

	case "move":
		targetWorkspace, all, err := parseMoveArgs(subArgs)
		if err != nil {
			fail(err)
		}

		targetIndex := targetWorkspace - 1
		_ = action.MoveToWorkspace(targetIndex, all)

	case "cycle":
		dir, err := parseCycleArgs(subArgs)
		if err != nil {
			fail(err)
		}

		_ = action.CycleWorkspace(dir)

	case "init":
		_ = action.InitWorkspaces()

	case "help", "-h", "--help", "":
		printUsage()

	default:
		fail(fmt.Errorf("unknown subcommand: %q", subcmd))
	}
}

func parseGotoArgs(args []string) (int, error) {
	if len(args) != 1 {
		return 0, errors.New("usage: hypr-local-workspaces goto <1..9>")
	}

	v, err := strconv.Atoi(args[0])
	if err != nil || v < 1 || v > 9 {
		return 0, errors.New("goto index must be a digit 1..9")
	}

	return v, nil
}

func parseMoveArgs(args []string) (int, bool, error) {
	fs := flag.NewFlagSet("move", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	all := fs.Bool("all", false, "Apply to all")

	// Separate flags from non-flag args
	if err := fs.Parse(args); err != nil {
		return 0, false, err
	}

	// After parsing, any leftover args are positional
	pos := fs.Args()
	if len(pos) != 1 {
		return 0, false, errors.New("usage: hypr-local-workspaces move <1..9> [--all]")
	}

	v, err := strconv.Atoi(pos[0])
	if err != nil {
		return 0, false, errors.New("move expects an integer digit")
	}

	if v < 1 || v > 9 {
		return 0, false, errors.New("move index must be a digit 1..9")
	}

	return v, *all, nil
}

func parseCycleArgs(args []string) (string, error) {
	if len(args) != 1 {
		return "", errors.New("usage: hypr-local-workspaces cycle <up|down>")
	}

	val := strings.ToLower(args[0])
	if val != "next" && val != "prev" {
		return "", errors.New("cycle direction must be 'next' or 'prev'")
	}

	return val, nil
}

func splitSubcommand(argv []string) (string, []string, error) {
	if len(argv) < 2 {
		return "", nil, nil
	}

	return argv[1], argv[2:], nil
}

func printUsage() {
	_, _ = fmt.Fprintln(os.Stderr, `Usage:
  hypr-local-workspaces goto  <1..9>
  hypr-local-workspaces move  <1..9> [--all]
  hypr-local-workspaces cycle <next|prev>`)
}

func fail(err error) {
	_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
	printUsage()
	os.Exit(ExitMissingArgs)
}
