package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MrMaksimize/granola-cli/internal/api"
	"github.com/MrMaksimize/granola-cli/internal/auth"
	"github.com/MrMaksimize/granola-cli/internal/format"
	"github.com/spf13/cobra"
)

var (
	showNotes      bool
	showTranscript bool
)

var showCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a meeting note",
	Long: `Show the content of a meeting note.

By default, shows the AI-generated summary. Use flags to select other content:
  --notes       Show your manually typed notes
  --transcript  Show the meeting transcript with timestamps`,
	Args: cobra.ExactArgs(1),
	RunE: runShow,
}

func init() {
	showCmd.Flags().BoolVar(&showNotes, "notes", false, "show your typed notes instead of AI summary")
	showCmd.Flags().BoolVar(&showTranscript, "transcript", false, "show meeting transcript")
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	docID := args[0]

	token, err := auth.LoadCredentials(credentialsPath)
	if err != nil {
		return err
	}

	client := newAPIClient(token)

	if showTranscript {
		return handleTranscript(client, docID)
	}

	if showNotes {
		return handleNotes(client, docID)
	}

	return handleSummary(client, docID)
}

func handleSummary(client *api.Client, docID string) error {
	doc, err := client.GetDocumentWithPanel(docID)
	if err != nil {
		return err
	}

	if doc.LastViewedPanel == nil || len(doc.LastViewedPanel) == 0 {
		fmt.Fprintln(os.Stderr, "No AI summary available for this note. Try --notes to see your typed notes.")
		return nil
	}

	var panel api.PanelContent
	if err := json.Unmarshal(doc.LastViewedPanel, &panel); err != nil {
		return fmt.Errorf("failed to parse panel data: %w", err)
	}

	md, err := format.ProseMirrorToMarkdown(panel.Content)
	if err != nil {
		// If ProseMirror parsing fails, the content might be a plain string
		var plainContent string
		if json.Unmarshal(panel.Content, &plainContent) == nil {
			md = plainContent
		} else {
			return fmt.Errorf("failed to convert summary to markdown: %w", err)
		}
	}

	if jsonOutput {
		return outputShowJSON(doc, md)
	}

	fmt.Printf("# %s\n", doc.Title)
	fmt.Printf("Created: %s\n\n", doc.CreatedAt.Format("Jan 02, 2006 3:04 PM"))
	fmt.Println(md)
	return nil
}

func handleNotes(client *api.Client, docID string) error {
	doc, err := client.GetDocument(docID)
	if err != nil {
		return err
	}

	if doc.NotesMarkdown == "" {
		fmt.Fprintln(os.Stderr, "No typed notes available for this note.")
		return nil
	}

	if jsonOutput {
		return outputShowJSON(doc, doc.NotesMarkdown)
	}

	fmt.Printf("# %s\n", doc.Title)
	fmt.Printf("Created: %s\n\n", doc.CreatedAt.Format("Jan 02, 2006 3:04 PM"))
	fmt.Println(doc.NotesMarkdown)
	return nil
}

func handleTranscript(client *api.Client, docID string) error {
	// First resolve the full document ID via prefix matching
	doc, err := client.GetDocument(docID)
	if err != nil {
		return err
	}

	utterances, err := client.GetTranscript(doc.ID)
	if err != nil {
		return err
	}

	if len(utterances) == 0 {
		fmt.Fprintln(os.Stderr, "No transcript available for this note.")
		return nil
	}

	if jsonOutput {
		return outputTranscriptJSON(doc, utterances)
	}

	fmt.Printf("# %s — Transcript\n", doc.Title)
	fmt.Printf("Created: %s\n\n", doc.CreatedAt.Format("Jan 02, 2006 3:04 PM"))

	for _, u := range utterances {
		source := u.Source
		if source == "microphone" {
			source = "You"
		} else if source == "system" {
			source = "Other"
		}
		timestamp := formatTimestamp(u.StartTimestamp)
		fmt.Printf("[%s] %s: %s\n", timestamp, source, u.Text)
	}
	return nil
}

func outputShowJSON(doc *api.Document, content string) error {
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
		Content:   content,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func outputTranscriptJSON(doc *api.Document, utterances []api.Utterance) error {
	out := struct {
		ID         string          `json:"id"`
		Title      string          `json:"title"`
		CreatedAt  string          `json:"created_at"`
		Utterances []api.Utterance `json:"utterances"`
	}{
		ID:         doc.ID,
		Title:      doc.Title,
		CreatedAt:  doc.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Utterances: utterances,
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}

func formatTimestamp(ts string) string {
	// Timestamps come as ISO strings; extract just HH:MM:SS
	if len(ts) >= 19 {
		parts := strings.Split(ts, "T")
		if len(parts) == 2 {
			timePart := parts[1]
			if idx := strings.Index(timePart, "."); idx > 0 {
				timePart = timePart[:idx]
			}
			if idx := strings.Index(timePart, "Z"); idx > 0 {
				timePart = timePart[:idx]
			}
			return timePart
		}
	}
	return ts
}
