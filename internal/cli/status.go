package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if proxy is running",
	Long:  "Check the status of the sandgrouse proxy.",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := readPID()
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "sandgrouse is not running")
			return
		}

		if !isProcessRunning(pid) {
			fmt.Fprintln(cmd.OutOrStdout(), "sandgrouse is not running (stale PID file cleaned up)")
			removePID()
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "sandgrouse is running (PID %d)\n", pid)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
