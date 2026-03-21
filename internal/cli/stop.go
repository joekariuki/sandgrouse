package cli

import (
	"fmt"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the proxy",
	Long:  "Stop the running sandgrouse proxy by sending SIGTERM.",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := readPID()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "sandgrouse is not running (no PID file found)")
			return
		}

		if !isProcessRunning(pid) {
			fmt.Fprintf(cmd.OutOrStdout(), "sandgrouse is not running (stale PID %d)\n", pid)
			removePID()
			return
		}

		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "failed to stop sandgrouse (PID %d): %v\n", pid, err)
			return
		}

		// Wait up to 5 seconds for process to exit
		for i := 0; i < 50; i++ {
			time.Sleep(100 * time.Millisecond)
			if !isProcessRunning(pid) {
				removePID()
				fmt.Fprintf(cmd.OutOrStdout(), "sandgrouse stopped (was PID %d)\n", pid)
				return
			}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "sandgrouse (PID %d) did not stop within 5 seconds\n", pid)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
