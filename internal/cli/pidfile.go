package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// sgDirOverride allows tests to redirect the sandgrouse directory.
var sgDirOverride string

// sandgrouseDir returns the path to ~/.sandgrouse/, creating it if needed.
func sandgrouseDir() (string, error) {
	if sgDirOverride != "" {
		if err := os.MkdirAll(sgDirOverride, 0o755); err != nil {
			return "", err
		}
		return sgDirOverride, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".sandgrouse")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// pidFilePath returns the path to the PID file.
func pidFilePath() (string, error) {
	dir, err := sandgrouseDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "sg.pid"), nil
}

// writePID writes a PID to the PID file.
func writePID(pid int) error {
	path, err := pidFilePath()
	if err != nil {
		return err
	}
	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0o644)
}

// readPID reads the PID from the PID file.
func readPID() (int, error) {
	path, err := pidFilePath()
	if err != nil {
		return 0, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

// removePID removes the PID file.
func removePID() error {
	path, err := pidFilePath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// isProcessRunning checks if a process with the given PID is alive.
func isProcessRunning(pid int) bool {
	err := syscall.Kill(pid, 0)
	return err == nil || errors.Is(err, syscall.EPERM)
}

// checkAlreadyRunning checks if sandgrouse is already running.
// Returns the PID and true if running, or 0 and false otherwise.
// Cleans up stale PID files automatically.
func checkAlreadyRunning() (int, bool) {
	pid, err := readPID()
	if err != nil {
		return 0, false
	}
	if isProcessRunning(pid) {
		return pid, true
	}
	// Stale PID file — process is dead, clean up
	removePID()
	return 0, false
}
