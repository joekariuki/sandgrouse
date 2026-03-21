package cli

import (
	"fmt"
	"io"

	"github.com/joekariuki/sandgrouse/internal/proxy"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "See bandwidth savings",
	Long:  "Show bandwidth savings from the sandgrouse proxy.",
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

		// Check if there's an active session (proxy running with session data)
		_, running := checkAlreadyRunning()
		if running && s.SessionRequests() > 0 {
			fmt.Fprintln(out, "Sandgrouse Stats")
			fmt.Fprintln(out, "────────────────")
			printStatsBlock(out, "This session",
				s.SessionRequests(),
				s.SessionRequestOriginalBytes(),
				s.SessionResponseOriginalBytes(),
				s.SessionResponseWireBytes())
			fmt.Fprintln(out)
			printStatsBlock(out, "All time",
				s.TotalRequests(),
				s.RequestOriginalBytes(),
				s.ResponseOriginalBytes(),
				s.ResponseWireBytes())
		} else {
			fmt.Fprintln(out, "Sandgrouse Stats")
			fmt.Fprintln(out, "────────────────")
			printStatsBlock(out, "All time",
				s.TotalRequests(),
				s.RequestOriginalBytes(),
				s.ResponseOriginalBytes(),
				s.ResponseWireBytes())
		}
	},
}

func printStatsBlock(out io.Writer, label string, reqs, reqBytes, respOrig, respWire int64) {
	fmt.Fprintf(out, "%s:\n", label)
	fmt.Fprintf(out, "  Requests:          %d\n", reqs)
	fmt.Fprintf(out, "  Request data:      %s\n", proxy.FormatBytes(reqBytes))
	if respOrig > 0 {
		saved := respOrig - respWire
		pct := float64(saved) / float64(respOrig) * 100
		fmt.Fprintf(out, "  Response savings:  %s (%.0f%% reduction)\n",
			proxy.FormatBytes(saved), pct)
	} else {
		fmt.Fprintln(out, "  Response savings:  no response data yet")
	}
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
