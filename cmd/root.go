package cmd

import (
	"cmd-ai-resolver/internal/logger"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	debug              bool
	passThroughCommand string
)

var (
	rootCmd = &cobra.Command{
		Use:   "cmd-ai-resolver [file-path]",
		Short: "Processes a shell command file with AI instructions.",
		Long: `cmd-ai-resolver is a CLI tool that takes a file path as an argument.
The file should contain a shell command, potentially with AI instruction tags
like <AI>your ai prompt</AI>. It processes these instructions using an LLM,
replaces the tags with the LLM's output, and saves the modified command back to the file.`,
		Args: cobra.ExactArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			logger.SetDebug(debug)
		},
		Run: func(cmd *cobra.Command, args []string) {
			filePath := args[0]
			cm := newCommandHandler(filePath)
			if err := cm.handleCommand(); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	Version = "v0.0.0+unknown"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print cmd-ai-resolver version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("cmd-ai-resolver version: %s\n", Version)
		},
	}
)

func init() {
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
	rootCmd.Flags().StringVarP(&passThroughCommand, "pass-through", "p", "", "Command to run if no AI tags are found")
	rootCmd.AddCommand(versionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
