package gocropus

import (
	"context"
	"os/exec"
)

// Cmd wraps information to run gocropus commands.
type Cmd struct {
	Exe   string // executable to run
	Model string // path of the model to use
}

// Run runs a command with the given arguments and returns its
// combined (stderr and stdout) output.
func (cmd *Cmd) Run(args ...string) ([]byte, error) {
	return cmd.RunContext(context.Background(), args...)
}

// RunContext runs the command with the given context.
func (cmd *Cmd) RunContext(ctx context.Context, args ...string) ([]byte, error) {
	args = cmd.Cmd(args...)
	runner := exec.CommandContext(ctx, args[0], args[1:]...)
	return runner.CombinedOutput()
}

// Cmd returns the command line command that get executed.
func (cmd *Cmd) Cmd(args ...string) []string {
	args = append([]string{cmd.Exe}, args...)
	if cmd.Model != "" {
		args = append(args, "--model", cmd.Model)
	}
	return args
}
