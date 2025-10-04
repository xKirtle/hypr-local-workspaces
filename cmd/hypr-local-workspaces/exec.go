package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

const (
	ExitSuccess     = 0   // Successful execution
	ExitFailure     = 1   // General failure
	ExitMissingArgs = 2   // Not enough arguments provided
	ExitTimeout     = 124 // Command timed out (like GNU timeout)
	ExitNotFound    = 127 // Command not found
	ExitInterrupted = 130 // Script terminated by Control-C
)

type CmdOptions struct {
	Stdin         io.Reader // nil = no stdin
	Stdout        io.Writer // nil + Capture=true -> internal buffer
	Stderr        io.Writer // nil + Capture=true && CombineOutput -> stdout buffer; else internal buffer
	Capture       bool      // capture stdout/stderr into buffers if no writer provided
	CombineOutput bool      // if Capture and Stderr nil -> redirect stderr into stdout buffer
	Interactive   bool      // inherit stdio from the process (ignored if Detached)
	DropStderr    bool      // send stderr to /dev/null
	Detached      bool      // start in new session; stdio redirected to /dev/null
	Timeout       time.Duration
	Dir           string
	Env           []string // nil -> inherit; non-nil -> replace
}

type CmdOpt func(*CmdOptions)

func WithInputString(s string) CmdOpt {
	return func(o *CmdOptions) { o.Stdin = bytes.NewBufferString(s) }
}

func WithInputBytes(b []byte) CmdOpt {
	return func(o *CmdOptions) {
		o.Stdin = bytes.NewReader(b)
	}
}
func CaptureOutput() CmdOpt {
	return func(o *CmdOptions) {
		o.Capture = true
	}
}
func CombineOutput() CmdOpt {
	return func(o *CmdOptions) {
		o.CombineOutput = true
	}
}
func Interactive() CmdOpt {
	return func(o *CmdOptions) {
		o.Interactive = true
	}
}
func DropStderr() CmdOpt {
	return func(o *CmdOptions) {
		o.DropStderr = true
	}
}
func Detached() CmdOpt {
	return func(o *CmdOptions) {
		o.Detached = true
	}
}
func WithTimeout(d time.Duration) CmdOpt {
	return func(o *CmdOptions) {
		o.Timeout = d
	}
}
func WithDir(dir string) CmdOpt {
	return func(o *CmdOptions) {
		o.Dir = dir
	}
}
func WithEnv(env []string) CmdOpt {
	return func(o *CmdOptions) {
		o.Env = env
	}
}
func WithStdout(w io.Writer) CmdOpt {
	return func(o *CmdOptions) {
		o.Stdout = w
	}
}
func WithStderr(w io.Writer) CmdOpt {
	return func(o *CmdOptions) {
		o.Stderr = w
	}
}

func Run(bin string, args ...string) ([]byte, int, error) {
	return RunWith(bin, args)
}

func RunWith(bin string, args []string, opts ...CmdOpt) ([]byte, int, error) {
	cmdOpts := CmdOptions{}
	for _, opt := range opts {
		opt(&cmdOpts)
	}

	if _, err := exec.LookPath(bin); err != nil {
		return nil, ExitNotFound, err
	}

	var (
		stdoutBuf bytes.Buffer
		stderrBuf bytes.Buffer
	)

	// Context (timeout)
	ctx := context.Background()
	var cancel context.CancelFunc
	if cmdOpts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), cmdOpts.Timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, bin, args...)
	if cmdOpts.Dir != "" {
		cmd.Dir = cmdOpts.Dir
	}
	if cmdOpts.Env != nil {
		cmd.Env = cmdOpts.Env
	}

	// Detached? redirect all FDs to /dev/null and start in new session
	if cmdOpts.Detached {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		null, _ := os.OpenFile(os.DevNull, os.O_RDWR, 0)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = null, null, null
		if err := cmd.Start(); err != nil {
			return nil, -1, err
		}

		return nil, 0, nil // started successfully
	}

	// Stdin
	cmd.Stdin = cmdOpts.Stdin

	// Stdout/Stderr wiring
	if cmdOpts.Interactive {
		cmd.Stdout = os.Stdout
		if cmdOpts.DropStderr {
			f, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
			defer func(f *os.File) {
				err := f.Close()
				if err != nil {

				}
			}(f)

			cmd.Stderr = f
		} else {
			cmd.Stderr = os.Stderr
		}
	} else if cmdOpts.Capture {
		// capture mode
		if cmdOpts.Stdout != nil {
			cmd.Stdout = cmdOpts.Stdout
		} else {
			cmd.Stdout = &stdoutBuf
		}

		if cmdOpts.Stderr != nil {
			cmd.Stderr = cmdOpts.Stderr
		} else if cmdOpts.CombineOutput {
			cmd.Stderr = cmd.Stdout
		} else {
			cmd.Stderr = &stderrBuf
		}
	} else {
		// default: inherit if writers not specified
		if cmdOpts.Stdout != nil {
			cmd.Stdout = cmdOpts.Stdout
		} else {
			cmd.Stdout = os.Stdout
		}

		if cmdOpts.Stderr != nil {
			cmd.Stderr = cmdOpts.Stderr
		} else {
			cmd.Stderr = os.Stderr
		}
	}

	err := cmd.Run()

	// Timeout?
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return stdoutBuf.Bytes(), ExitFailure, ctx.Err()
	}

	if err == nil {
		// success
		if cmdOpts.Capture {
			// If not combining, append stderr to stdout like your previous helpers
			if !cmdOpts.CombineOutput && stderrBuf.Len() > 0 {
				stdoutBuf.Write(stderrBuf.Bytes())
			}

			return stdoutBuf.Bytes(), 0, nil
		}

		return nil, 0, nil
	}

	var ee *exec.ExitError
	if errors.As(err, &ee) {
		if cmdOpts.Capture {
			if !cmdOpts.CombineOutput && stderrBuf.Len() > 0 {
				stdoutBuf.Write(stderrBuf.Bytes())
			}

			return stdoutBuf.Bytes(), ee.ExitCode(), nil
		}

		return nil, ee.ExitCode(), nil
	}

	return nil, -1, err
}

func MustRun(bin string, args ...string) ([]byte, int) {
	out, code, err := Run(bin, args...)
	if err != nil {
		log.Fatalf("running %s %v: %v", bin, args, err)
	}

	return out, code
}

func MustRunWith(bin string, args []string, opts ...CmdOpt) ([]byte, int) {
	out, code, err := RunWith(bin, args, opts...)
	if err != nil {
		log.Fatalf("running %s %v: %v", bin, args, err)
	}

	return out, code
}
