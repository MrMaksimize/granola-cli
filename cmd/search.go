package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/MrMaksimize/granola-cli/internal/api"
	"github.com/MrMaksimize/granola-cli/internal/auth"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search meeting notes",
	Long:  `Search meeting notes by title (case-insensitive).`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := strings.ToLower(args[0])

	token, err := auth.LoadCredentials(credentialsPath)
	if err != nil {
		return err
	}

	client := api.NewClient(token)
	docs, err := client.ListDocuments(100, 0)
	if err != nil {
		return err
	}

	// Client-side filtering
	var matches []api.Document
	for _, doc := range docs {
		if strings.Contains(strings.ToLower(doc.Title), query) {
			matches = append(matches, doc)
		}
	}

	if jsonOutput {
		return outputSearchJSON(matches)
	}

	return outputSearchText(matches, query)
}

func outputSearchJSON(docs []api.Document) error {
	type jsonDoc struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		CreatedAt string `json:"created_at"`
	}

	out := make([]jsonDoc, len(docs))
	for i, doc := range docs {
		out[i] = jsonDoc{
			ID:        doc.ID,
			Title:     doc.Title,
			CreatedAt: doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func outputSearchText(docs []api.Document, query string) error {
	if len(docs) == 0 {
		fmt.Printf("No notes found matching '%s'\n", query)
		return nil
	}

	fmt.Printf("Found %d note(s) matching '%s'\n\n", len(docs), query)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tCreated\tTitle")
	fmt.Fprintln(w, "──\t───────\t─────")

	for _, doc := range docs {
		title := doc.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n",
			doc.ID[:8],
			doc.CreatedAt.Format("Jan 02"),
			title,
		)
	}

	return w.Flush()
}
