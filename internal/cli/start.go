package cli

import (
	"fmt"
	"log"

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

		fmt.Fprint(cmd.OutOrStdout(), banner)
		fmt.Fprintf(cmd.OutOrStdout(), "v%s | Stop burning data bundles on AI tools.\n\n", Version)

		srv := &proxy.Server{
			ListenAddr: addr,
			Algorithm:  algo,
		}
		if err := srv.Start(); err != nil {
			log.Fatalf("failed to start proxy: %v", err)
		}
	},
}

func init() {
	startCmd.Flags().StringP("addr", "a", ":8080", "address to listen on")
	startCmd.Flags().String("algorithm", "brotli", "compression algorithm (gzip or brotli)")
	rootCmd.AddCommand(startCmd)
}
