package api

import "time"

// Document represents a Granola meeting note
type Document struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	NotesMarkdown string    `json:"notes_markdown"`
	NotesPlain    string    `json:"notes_plain"`
	Summary       string    `json:"summary"`
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
