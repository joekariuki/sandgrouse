package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is the current sandgrouse version.
// Overridden at build time via ldflags for releases.
var Version = "0.1.0-dev"

const banner = " ____                  _\n" +
	"/ ___|  __ _ _ __   __| | __ _ _ __ ___  _   _ ___  ___\n" +
	"\\___ \\ / _` | '_ \\ / _` |/ _` | '__/ _ \\| | | / __|/ _ \\\n" +
	" ___) | (_| | | | | (_| | (_| | | | (_) | |_| \\__ \\  __/\n" +
	"|____/ \\__,_|_| |_|\\__,_|\\__, |_|  \\___/ \\__,_|___/\\___|\n" +
	"                          |___/\n"

var welcome = banner +
	"v%s | Stop burning data bundles on AI tools.\n" +
	"\n" +
	"sandgrouse compresses LLM API traffic for developers on metered connections.\n" +
	"\n" +
	"Quick start:\n" +
	"  sg start                  Start the proxy\n" +
	"  sg status                 Check if proxy is running\n" +
	"  sg stats                  See bandwidth savings\n" +
	"  sg stop                   Stop the proxy\n" +
	"\n" +
	"Setup your AI tools:\n" +
	"  export ANTHROPIC_BASE_URL=http://localhost:8080   # Claude Code\n" +
	"  export OPENAI_BASE_URL=http://localhost:8080       # Cursor / OpenAI\n"

var rootCmd = &cobra.Command{
	Use:   "sg",
	Short: "Sandgrouse - LLM traffic compression proxy",
	Long:  "Sandgrouse compresses LLM API traffic for developers on metered connections.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintf(cmd.OutOrStdout(), welcome, Version)
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
