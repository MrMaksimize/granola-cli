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
	Long:  `Search meeting notes by title and content (case-insensitive).`,
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

	client := newAPIClient(token)
	docs, err := client.ListDocuments(100, 0)
	if err != nil {
		return err
	}

	// Client-side filtering: match title and content
	var matches []searchMatch
	for _, doc := range docs {
		inTitle := strings.Contains(strings.ToLower(doc.Title), query)
		inContent := strings.Contains(strings.ToLower(doc.NotesMarkdown), query)
		if inTitle || inContent {
			matches = append(matches, searchMatch{Doc: doc, InTitle: inTitle, InContent: inContent})
		}
	}

	if jsonOutput {
		matchDocs := make([]api.Document, len(matches))
		for i, m := range matches {
			matchDocs[i] = m.Doc
		}
		return outputSearchJSON(matchDocs)
	}

	return outputSearchTextWithLocation(matches, query)
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

type searchMatch struct {
	Doc       api.Document
	InTitle   bool
	InContent bool
}

func outputSearchTextWithLocation(matches []searchMatch, query string) error {
	if len(matches) == 0 {
		fmt.Printf("No notes found matching '%s'\n", query)
		return nil
	}

	fmt.Printf("Found %d note(s) matching '%s'\n\n", len(matches), query)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tCreated\tMatch\tTitle")
	fmt.Fprintln(w, "──\t───────\t─────\t─────")

	for _, m := range matches {
		title := m.Doc.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}
		matchLoc := matchLocation(m.InTitle, m.InContent)
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
			m.Doc.ID[:8],
			m.Doc.CreatedAt.Format("Jan 02"),
			matchLoc,
			title,
		)
	}

	return w.Flush()
}

func matchLocation(inTitle, inContent bool) string {
	if inTitle && inContent {
		return "title+content"
	}
	if inTitle {
		return "title"
	}
	return "content"
}
