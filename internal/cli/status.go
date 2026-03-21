package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if proxy is running",
	Long:  "Check the status of the sandgrouse proxy. Requires daemon mode (coming in a future release).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "sg status is not implemented yet — daemon mode coming soon.")
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
