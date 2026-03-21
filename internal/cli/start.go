package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/joekariuki/sandgrouse/internal/proxy"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the proxy",
	Long:  "Start the sandgrouse compression proxy. Traffic is proxied to upstream LLM APIs with response compression.",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Flags().GetString("addr")
		algo, _ := cmd.Flags().GetString("algorithm")
		foreground, _ := cmd.Flags().GetBool("foreground")

		if foreground {
			runForeground(cmd, addr, algo)
		} else {
			runDaemon(cmd, addr, algo)
		}
	},
}

func init() {
	startCmd.Flags().StringP("addr", "a", ":8080", "address to listen on")
	startCmd.Flags().String("algorithm", "brotli", "compression algorithm (gzip or brotli)")
	startCmd.Flags().BoolP("foreground", "f", false, "run in foreground (default: daemon mode)")
	rootCmd.AddCommand(startCmd)
}

// statsFilePath returns the path to the stats JSON file.
func statsFilePath() string {
	dir, err := sandgrouseDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "stats.json")
}

// runForeground starts the proxy in the current process with signal handling.
func runForeground(cmd *cobra.Command, addr, algo string) {
	// Write our PID so sg stop/status can find us
	if err := writePID(os.Getpid()); err != nil {
		log.Printf("warning: could not write PID file: %v", err)
	}
	defer removePID()

	fmt.Fprint(cmd.OutOrStdout(), banner)
	fmt.Fprintf(cmd.OutOrStdout(), "v%s | Stop burning data bundles on AI tools.\n\n", Version)

	// Load persisted stats from previous sessions
	stats := &proxy.Stats{}
	statsPath := statsFilePath()
	if statsPath != "" {
		if err := stats.LoadFrom(statsPath); err != nil {
			log.Printf("warning: could not load stats: %v", err)
		}
	}

	srv := &proxy.Server{
		ListenAddr: addr,
		Algorithm:  algo,
	}
	srv.SetStats(stats)

	// Start server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Save stats periodically (every 60 seconds)
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			if statsPath != "" {
				if err := stats.SaveTo(statsPath); err != nil {
					log.Printf("warning: could not save stats: %v", err)
				}
			}
		}
	}()

	// Wait for signal or server error
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("received %s, shutting down...", sig)
		// Save stats before shutdown
		if statsPath != "" {
			if err := stats.SaveTo(statsPath); err != nil {
				log.Printf("warning: could not save stats on shutdown: %v", err)
			}
		}
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
		log.Println("sandgrouse stopped")
	case err := <-errCh:
		if err != nil {
			log.Fatalf("failed to start proxy: %v", err)
		}
	}
}

// runDaemon spawns the proxy as a detached background process.
func runDaemon(cmd *cobra.Command, addr, algo string) {
	// Check if already running
	if pid, running := checkAlreadyRunning(); running {
		fmt.Fprintf(cmd.OutOrStdout(), "sandgrouse is already running (PID %d)\n", pid)
		return
	}

	// Find our own executable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("could not find executable path: %v", err)
	}

	// Spawn child in foreground mode, detached from this terminal
	child := exec.Command(exePath, "start", "--foreground", "--addr", addr, "--algorithm", algo)
	child.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	child.Stdout = nil
	child.Stderr = nil

	if err := child.Start(); err != nil {
		log.Fatalf("failed to start daemon: %v", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "sandgrouse started (PID %d)\n", child.Process.Pid)
	fmt.Fprintf(cmd.OutOrStdout(), "Proxy listening on %s\n", addr)
	fmt.Fprintf(cmd.OutOrStdout(), "Stop with: sg stop\n")
}
