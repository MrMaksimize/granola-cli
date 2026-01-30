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

var foldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List document folders",
	Long:  `List all document folders/lists from Granola.`,
	RunE:  runFolders,
}

func init() {
	rootCmd.AddCommand(foldersCmd)
}

func runFolders(cmd *cobra.Command, args []string) error {
	token, err := auth.LoadCredentials(credentialsPath)
	if err != nil {
		return err
	}

	client := newAPIClient(token)
	lists, err := client.GetDocumentLists()
	if err != nil {
		return err
	}

	if jsonOutput {
		return outputFoldersJSON(lists)
	}

	return outputFoldersText(lists)
}

func outputFoldersJSON(lists []api.DocumentList) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(lists)
}

func outputFoldersText(lists []api.DocumentList) error {
	if len(lists) == 0 {
		fmt.Println("No folders found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "Name\tDocs\tCreated")
	fmt.Fprintln(w, "────\t────\t───────")

	for _, l := range lists {
		created := l.CreatedAt
		if len(created) >= 10 {
			created = created[:10]
		}
		name := l.Name
		if len(name) > 40 {
			name = name[:37] + "..."
		}
		favourite := ""
		if l.IsFavourite {
			favourite = " *"
		}
		fmt.Fprintf(w, "%s%s\t%d\t%s\n",
			name,
			favourite,
			l.DocumentCount,
			created,
		)
	}

	return w.Flush()
}
