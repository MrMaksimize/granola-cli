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
		ID              string                 `json:"id"`
		Title           string                 `json:"title"`
		CreatedAt       string                 `json:"created_at"`
		UpdatedAt       string                 `json:"updated_at"`
		Content         string                 `json:"content"`
		Summary         string                 `json:"summary,omitempty"`
		LastViewedPanel map[string]interface{} `json:"last_viewed_panel,omitempty"`
	}{
		ID:              doc.ID,
		Title:           doc.Title,
		CreatedAt:       doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       doc.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		Content:         doc.NotesMarkdown,
		Summary:         doc.Summary,
		LastViewedPanel: doc.LastViewedPanel,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func outputShowText(doc *api.Document) error {
	fmt.Printf("# %s\n", doc.Title)
	fmt.Printf("Created: %s\n\n", doc.CreatedAt.Format("Jan 02, 2006 3:04 PM"))

	// Try to extract enhanced notes from last_viewed_panel
	if doc.LastViewedPanel != nil {
		if content, ok := doc.LastViewedPanel["content"]; ok {
			if enhanced := extractProseMirrorText(content, 0); enhanced != "" {
				fmt.Println(enhanced)
				return nil
			}
		}
	}

	// Fall back to basic notes
	if doc.Summary != "" {
		fmt.Println("## Summary")
		fmt.Println(doc.Summary)
		fmt.Println()
		fmt.Println("## Notes")
	}
	fmt.Println(doc.NotesMarkdown)
	return nil
}

// extractProseMirrorText recursively extracts text from ProseMirror JSON content
func extractProseMirrorText(node interface{}, depth int) string {
	switch n := node.(type) {
	case map[string]interface{}:
		nodeType, _ := n["type"].(string)
		content, hasContent := n["content"]

		switch nodeType {
		case "text":
			if text, ok := n["text"].(string); ok {
				return text
			}
		case "heading":
			level := 3
			if attrs, ok := n["attrs"].(map[string]interface{}); ok {
				if l, ok := attrs["level"].(float64); ok {
					level = int(l)
				}
			}
			prefix := ""
			for i := 0; i < level; i++ {
				prefix += "#"
			}
			if hasContent {
				return "\n" + prefix + " " + extractProseMirrorText(content, depth) + "\n"
			}
		case "paragraph":
			if hasContent {
				return extractProseMirrorText(content, depth) + "\n"
			}
			return "\n"
		case "bulletList", "orderedList":
			if hasContent {
				return extractProseMirrorText(content, depth)
			}
		case "listItem":
			if hasContent {
				return "- " + extractProseMirrorText(content, depth+1)
			}
		default:
			if hasContent {
				return extractProseMirrorText(content, depth)
			}
		}
	case []interface{}:
		var result string
		for _, item := range n {
			result += extractProseMirrorText(item, depth)
		}
		return result
	}
	return ""
}
