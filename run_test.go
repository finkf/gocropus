package gocropus

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestCmdArgs(t *testing.T) {
	tests := []struct {
		cmd  *Cmd
		args []string
		want string
	}{
		{&Cmd{Exe: "test"}, nil, "test"},
		{&Cmd{Exe: "test", Model: "test"}, nil, "test --model test"},
		{&Cmd{Exe: "test", Model: "test"}, []string{"a", "b"}, "test a b --model test"},
		{&Cmd{Exe: "test"}, []string{"a", "b"}, "test a b"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			if got := strings.Join(tc.cmd.Cmd(tc.args...), " "); got != tc.want {
				t.Fatalf("expected %q; got %q", tc.want, got)
			}
		})
	}
}

func TestCmdOutput(t *testing.T) {
	cmd := &Cmd{Exe: "testdata/print-model.sh", Model: "testmodel"}
	out, err := cmd.Run("a", "b", "c")
	if err != nil {
		t.Fatalf("got error: %v", err)
	}
	if got := strings.TrimSpace(string(out)); got != "testmodel" {
		t.Fatalf("expected %q; got %q", "testmodel", got)
	}
}

func TestCmdTimeout(t *testing.T) {
	cmd := &Cmd{Exe: "testdata/sleep.sh"}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := cmd.RunContext(ctx, "a", "b", "c")
	if err == nil {
		t.Fatalf("expected an error")
	}
}
