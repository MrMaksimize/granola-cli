package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/MrMaksimize/granola-cli/internal/api"
	"github.com/MrMaksimize/granola-cli/internal/auth"
	"github.com/spf13/cobra"
)

var listLimit int

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent meeting notes",
	Long:  `List recent meeting notes from Granola.`,
	RunE:  runList,
}

func init() {
	listCmd.Flags().IntVar(&listLimit, "limit", 20, "number of notes to fetch")
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	token, err := auth.LoadCredentials(credentialsPath)
	if err != nil {
		return err
	}

	client := newAPIClient(token)
	docs, err := client.ListDocuments(listLimit, 0)
	if err != nil {
		return err
	}

	if jsonOutput {
		return outputListJSON(docs)
	}

	return outputListText(docs)
}

func outputListJSON(docs []api.Document) error {
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

func outputListText(docs []api.Document) error {
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
