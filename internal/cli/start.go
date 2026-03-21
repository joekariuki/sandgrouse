package cli

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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

		fmt.Fprint(cmd.OutOrStdout(), banner)
		fmt.Fprintf(cmd.OutOrStdout(), "v%s | Stop burning data bundles on AI tools.\n\n", Version)

		srv := &proxy.Server{
			ListenAddr: addr,
			Algorithm:  algo,
		}

		// Start server in a goroutine
		errCh := make(chan error, 1)
		go func() {
			errCh <- srv.Start()
		}()

		// Wait for signal or server error
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-sigCh:
			log.Printf("received %s, shutting down...", sig)
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
	},
}

func init() {
	startCmd.Flags().StringP("addr", "a", ":8080", "address to listen on")
	startCmd.Flags().String("algorithm", "brotli", "compression algorithm (gzip or brotli)")
	rootCmd.AddCommand(startCmd)
}
