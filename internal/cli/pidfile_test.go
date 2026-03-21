package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func withTempDir(t *testing.T) {
	t.Helper()
	old := sgDirOverride
	sgDirOverride = filepath.Join(t.TempDir(), ".sandgrouse")
	t.Cleanup(func() { sgDirOverride = old })
}

func TestWriteAndReadPID(t *testing.T) {
	withTempDir(t)

	if err := writePID(12345); err != nil {
		t.Fatalf("writePID() error: %v", err)
	}

	pid, err := readPID()
	if err != nil {
		t.Fatalf("readPID() error: %v", err)
	}
	if pid != 12345 {
		t.Errorf("readPID() = %d, want 12345", pid)
	}
}

func TestReadPIDFileNotExist(t *testing.T) {
	withTempDir(t)

	_, err := readPID()
	if err == nil {
		t.Fatal("readPID() expected error for missing file")
	}
}

func TestReadPIDInvalidContent(t *testing.T) {
	withTempDir(t)

	path, _ := pidFilePath()
	os.MkdirAll(filepath.Dir(path), 0o755)
	os.WriteFile(path, []byte("notanumber"), 0o644)

	_, err := readPID()
	if err == nil {
		t.Fatal("readPID() expected error for invalid content")
	}
}

func TestRemovePID(t *testing.T) {
	withTempDir(t)

	writePID(12345)
	if err := removePID(); err != nil {
		t.Fatalf("removePID() error: %v", err)
	}

	_, err := readPID()
	if err == nil {
		t.Fatal("PID file should not exist after removePID()")
	}
}

func TestRemovePIDNoFile(t *testing.T) {
	withTempDir(t)

	if err := removePID(); err != nil {
		t.Fatalf("removePID() should not error when file missing: %v", err)
	}
}

func TestIsProcessRunning(t *testing.T) {
	// Current process should be running
	if !isProcessRunning(os.Getpid()) {
		t.Error("isProcessRunning(os.Getpid()) = false, want true")
	}

	// Non-existent process should not be running
	if isProcessRunning(99999999) {
		t.Error("isProcessRunning(99999999) = true, want false")
	}
}

func TestCheckAlreadyRunningStaleCleanup(t *testing.T) {
	withTempDir(t)

	// Write a PID for a dead process
	writePID(99999999)

	pid, running := checkAlreadyRunning()
	if running {
		t.Errorf("checkAlreadyRunning() = (%d, true), want (0, false)", pid)
	}

	// Stale PID file should have been cleaned up
	_, err := readPID()
	if err == nil {
		t.Fatal("stale PID file should have been removed")
	}
}

func TestCheckAlreadyRunningLiveProcess(t *testing.T) {
	withTempDir(t)

	// Write our own PID — we're definitely alive
	writePID(os.Getpid())

	pid, running := checkAlreadyRunning()
	if !running {
		t.Fatal("checkAlreadyRunning() = (_, false), want true")
	}
	if pid != os.Getpid() {
		t.Errorf("checkAlreadyRunning() pid = %d, want %d", pid, os.Getpid())
	}
}

func TestSandgrouseDirCreation(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "nested", ".sandgrouse")
	sgDirOverride = newDir
	t.Cleanup(func() { sgDirOverride = "" })

	dir, err := sandgrouseDir()
	if err != nil {
		t.Fatalf("sandgrouseDir() error: %v", err)
	}
	if dir != newDir {
		t.Errorf("sandgrouseDir() = %q, want %q", dir, newDir)
	}

	info, err := os.Stat(newDir)
	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("sandgrouseDir() did not create a directory")
	}
}
