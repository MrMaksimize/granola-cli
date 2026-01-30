package cmd

import (
	"os"

	"github.com/MrMaksimize/granola-cli/internal/api"
	"github.com/spf13/cobra"
)

var (
	credentialsPath string
	jsonOutput      bool
	apiURL          string
)

var rootCmd = &cobra.Command{
	Use:   "granola",
	Short: "CLI for Granola meeting notes",
	Long:  `A command-line interface for interacting with Granola meeting notes.`,
}

// Execute runs the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&credentialsPath, "credentials", "", "path to credentials file (default: ~/Library/Application Support/Granola/supabase.json)")
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output as JSON")
	rootCmd.PersistentFlags().StringVar(&apiURL, "api-url", "", "override API base URL (for testing)")
	rootCmd.PersistentFlags().MarkHidden("api-url")
}

// newAPIClient creates an API client, using the override URL if set
func newAPIClient(token string) *api.Client {
	if apiURL != "" {
		return api.NewClientWithBaseURL(token, apiURL)
	}
	return api.NewClient(token)
}
