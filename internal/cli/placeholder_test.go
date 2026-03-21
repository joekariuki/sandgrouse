package cli

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestCommandOutputs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantMsg string
	}{
		{
			name:    "stop shows not running when no PID",
			args:    []string{"stop"},
			wantMsg: "not running",
		},
		{
			name:    "status shows not running when no PID",
			args:    []string{"status"},
			wantMsg: "not running",
		},
		{
			name:    "stats shows no data when empty",
			args:    []string{"stats"},
			wantMsg: "no stats recorded yet",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer resetRootCmd()
			// Use temp dir so no real PID file is found
			old := sgDirOverride
			sgDirOverride = filepath.Join(t.TempDir(), ".sandgrouse")
			defer func() { sgDirOverride = old }()

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("Execute(%v) error: %v", tt.args, err)
			}

			if !strings.Contains(buf.String(), tt.wantMsg) {
				t.Errorf("output missing %q\ngot: %s", tt.wantMsg, buf.String())
			}
		})
	}
}

func TestAllCommandsRegistered(t *testing.T) {
	wantCmds := []string{"start", "stop", "status", "stats"}

	for _, name := range wantCmds {
		cmd, _, err := rootCmd.Find([]string{name})
		if err != nil {
			t.Errorf("command %q not found: %v", name, err)
			continue
		}
		if cmd.Use != name {
			t.Errorf("Find(%q) returned command with Use=%q", name, cmd.Use)
		}
	}
}

func TestHelpShowsAllCommands(t *testing.T) {
	defer resetRootCmd()
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	rootCmd.Execute()

	output := buf.String()
	for _, name := range []string{"start", "stop", "status", "stats"} {
		if !strings.Contains(output, name) {
			t.Errorf("help output missing command %q", name)
		}
	}
}
