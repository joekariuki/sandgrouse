package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "See bandwidth savings",
	Long:  "Show bandwidth savings from the sandgrouse proxy. Requires stats persistence (coming in a future release).",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "sg stats is not implemented yet — stats persistence coming soon.")
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
