package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCommandShowsBanner(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})
	defer rootCmd.SetOut(nil)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error: %v", err)
	}

	output := buf.String()
	wantParts := []string{
		"sandgrouse",
		Version,
		"Stop burning data bundles",
		"sg start",
		"sg stop",
		"ANTHROPIC_BASE_URL",
	}
	for _, part := range wantParts {
		if !strings.Contains(output, part) {
			t.Errorf("output missing %q", part)
		}
	}
}

func TestStartCommandRegistered(t *testing.T) {
	cmd, _, err := rootCmd.Find([]string{"start"})
	if err != nil {
		t.Fatalf("Find('start') error: %v", err)
	}
	if cmd.Use != "start" {
		t.Errorf("Use = %q, want 'start'", cmd.Use)
	}
}

func TestStartCommandHasFlags(t *testing.T) {
	cmd, _, _ := rootCmd.Find([]string{"start"})

	addr := cmd.Flags().Lookup("addr")
	if addr == nil {
		t.Fatal("start command missing --addr flag")
	}
	if addr.DefValue != ":8080" {
		t.Errorf("--addr default = %q, want ':8080'", addr.DefValue)
	}

	algo := cmd.Flags().Lookup("algorithm")
	if algo == nil {
		t.Fatal("start command missing --algorithm flag")
	}
	if algo.DefValue != "brotli" {
		t.Errorf("--algorithm default = %q, want 'brotli'", algo.DefValue)
	}
}
