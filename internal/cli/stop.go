package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the proxy",
	Long:  "Stop the sandgrouse proxy. Requires daemon mode (coming in a future release).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "sg stop is not implemented yet — daemon mode coming soon.")
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
