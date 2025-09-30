package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const (
	ExitSuccess     = 0   // Successful execution
	ExitFailure     = 1   // General failure
	ExitMissingArgs = 2   // Not enough arguments provided
	ExitTimeout     = 124 // Command timed out (like GNU timeout)
	ExitNotFound    = 127 // Command not found
	ExitInterrupted = 130 // Script terminated by Control-C
)

// Fatalf prints a message and exits with the given code
func Fatalf(code int, format string, args ...any) {
	log.SetFlags(0)
	log.Printf(format, args...)
	os.Exit(code)
}

// Check maps common error kinds to shell-like exit codes and exits.
func Check(err error, context string, args ...any) {
	if err == nil {
		return
	}
	msg := fmt.Sprintf(context, args...)
	Fatalf(classifyExitCode(err), "%s: %v", msg, err)
}

// classifyExitCode centralizes the policy for exit codes.
func classifyExitCode(err error) int {
	switch {
	case errors.Is(err, exec.ErrNotFound):
		return ExitNotFound
	case errors.Is(err, context.DeadlineExceeded):
		return ExitTimeout
	case errors.Is(err, context.Canceled):
		return ExitInterrupted
	default:
		return ExitFailure
	}
}
