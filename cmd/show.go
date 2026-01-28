package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MrMaksimize/granola-cli/internal/api"
	"github.com/MrMaksimize/granola-cli/internal/auth"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a meeting note",
	Long:  `Show the full content of a meeting note.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	docID := args[0]

	token, err := auth.LoadCredentials(credentialsPath)
	if err != nil {
		return err
	}

	client := api.NewClient(token)
	doc, err := client.GetDocument(docID)
	if err != nil {
		return err
	}

	if jsonOutput {
		return outputShowJSON(doc)
	}

	return outputShowText(doc)
}

func outputShowJSON(doc *api.Document) error {
	out := struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Content   string `json:"content"`
	}{
		ID:        doc.ID,
		Title:     doc.Title,
		CreatedAt: doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Content:   doc.NotesMarkdown,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func outputShowText(doc *api.Document) error {
	fmt.Printf("# %s\n", doc.Title)
	fmt.Printf("Created: %s\n\n", doc.CreatedAt.Format("Jan 02, 2006 3:04 PM"))
	fmt.Println(doc.NotesMarkdown)
	return nil
}
