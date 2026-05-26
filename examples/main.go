package main

import (
	"fmt"
	"os"

	wizard "github.com/sorrowcode/goactive-cli"
	"github.com/spf13/cobra"
)

// 1. They import your open-source library
// ==========================================
// DEVELOPER's STANDARD COBRA CODE
// ==========================================

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "A CLI for managing a local server and sending messages",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage the local server",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		detached, _ := cmd.Flags().GetBool("detached")

		fmt.Printf("🟢 Starting server on port %d...\n", port)
		if detached {
			fmt.Println("👻 Running in background (detached mode).")
		}
	},
}

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send a message payload",
	Run: func(cmd *cobra.Command, args []string) {
		payload, _ := cmd.Flags().GetString("payload")
		fmt.Printf("📬 Sending Payload:\n%s\n", payload)
	},
}

func init() {
	// Configure server start flags
	startCmd.Flags().Int("port", 8080, "Port to run the server on")
	startCmd.Flags().Bool("detached", false, "Run server in the background")

	// Configure send payload flags
	defaultPayload := `{"hello": "world", "status": "testing"}`
	sendCmd.Flags().String("payload", defaultPayload, "JSON Payload to send")

	// They can use your multiline annotation if they want the large text editor!
	sendCmd.Flags().SetAnnotation("payload", "editor", []string{"multiline"})

	// Build the command tree
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(startCmd)
	rootCmd.AddCommand(sendCmd)
}

// ==========================================
// THE MAGIC INTEGRATION
// ==========================================

func main() {
	// 2. Instead of calling rootCmd.Execute(), they wrap it in your library.
	// If the user runs `myapp server start`, it executes normally.
	// If the user runs `myapp`, your interactive wizard takes over!
	if err := wizard.Execute(rootCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
