package api

import (
	"encoding/json"
	"time"
)

// Document represents a Granola meeting note
type Document struct {
	ID              string          `json:"id"`
	Title           string          `json:"title"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	NotesMarkdown   string          `json:"notes_markdown"`
	NotesPlain      string          `json:"notes_plain"`
	Summary         string          `json:"summary"`
	LastViewedPanel json.RawMessage `json:"last_viewed_panel,omitempty"`
}

// ListDocumentsRequest is the request body for get-documents
type ListDocumentsRequest struct {
	Limit                  int  `json:"limit"`
	Offset                 int  `json:"offset"`
	IncludeLastViewedPanel bool `json:"include_last_viewed_panel"`
}

// ListDocumentsResponse is the response from get-documents
type ListDocumentsResponse struct {
	Docs []Document `json:"docs"`
}

// PanelContent represents the last_viewed_panel structure
type PanelContent struct {
	Content json.RawMessage `json:"content"`
}

// Utterance represents a single transcript utterance
type Utterance struct {
	ID             string `json:"id"`
	DocumentID     string `json:"document_id"`
	Source         string `json:"source"`
	Text           string `json:"text"`
	StartTimestamp string `json:"start_timestamp"`
	EndTimestamp   string `json:"end_timestamp"`
	IsFinal        bool   `json:"is_final"`
}

// TranscriptResponse is the response from get-document-transcript
type TranscriptResponse struct {
	Utterances []Utterance `json:"utterances"`
}

// DocumentList represents a folder/list of documents
type DocumentList struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CreatedAt     string `json:"created_at"`
	WorkspaceID   string `json:"workspace_id"`
	DocumentCount int    `json:"document_count"`
	IsFavourite   bool   `json:"is_favourite"`
}

// DocumentListsResponse is the response from get-document-lists
type DocumentListsResponse struct {
	Lists []DocumentList `json:"lists"`
}
