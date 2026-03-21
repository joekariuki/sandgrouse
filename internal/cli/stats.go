package cli

import (
	"fmt"

	"github.com/joekariuki/sandgrouse/internal/proxy"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "See bandwidth savings",
	Long:  "Show cumulative bandwidth savings from the sandgrouse proxy.",
	Run: func(cmd *cobra.Command, args []string) {
		statsPath := statsFilePath()
		if statsPath == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "could not determine stats file path")
			return
		}

		s := &proxy.Stats{}
		if err := s.LoadFrom(statsPath); err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "no stats recorded yet")
			return
		}

		if s.TotalRequests() == 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "no stats recorded yet")
			return
		}

		out := cmd.OutOrStdout()
		fmt.Fprintln(out, "Sandgrouse Stats")
		fmt.Fprintln(out, "────────────────")
		fmt.Fprintf(out, "Requests proxied:    %d\n", s.TotalRequests())
		fmt.Fprintf(out, "Request data:        %s\n", proxy.FormatBytes(s.RequestOriginalBytes()))

		respOrig := s.ResponseOriginalBytes()
		respWire := s.ResponseWireBytes()
		fmt.Fprintf(out, "Response data:       %s original, %s on wire\n",
			proxy.FormatBytes(respOrig), proxy.FormatBytes(respWire))

		if respOrig > 0 {
			saved := respOrig - respWire
			pct := float64(saved) / float64(respOrig) * 100
			fmt.Fprintf(out, "Response savings:    %s (%.0f%% reduction)\n",
				proxy.FormatBytes(saved), pct)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
