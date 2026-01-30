package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	defaultBaseURL = "https://api.granola.ai"
	clientVersion  = "5.354.0"
)

// Client is a Granola API client
type Client struct {
	token      string
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Granola API client
func NewClient(token string) *Client {
	return &Client{
		token:      token,
		baseURL:    defaultBaseURL,
		httpClient: &http.Client{},
	}
}

// NewClientWithBaseURL creates a client with a custom base URL (for testing)
func NewClientWithBaseURL(token, baseURL string) *Client {
	return &Client{
		token:      token,
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

// ListDocuments fetches documents from the Granola API
func (c *Client) ListDocuments(limit, offset int) ([]Document, error) {
	req := ListDocumentsRequest{
		Limit:                  limit,
		Offset:                 offset,
		IncludeLastViewedPanel: false,
	}
	return c.fetchDocuments(req)
}

// GetDocument fetches a single document by ID (prefix matching supported)
func (c *Client) GetDocument(id string) (*Document, error) {
	docs, err := c.ListDocuments(100, 0)
	if err != nil {
		return nil, err
	}
	return c.findDocByID(docs, id)
}

// GetDocumentWithPanel fetches a document with the last_viewed_panel included
func (c *Client) GetDocumentWithPanel(id string) (*Document, error) {
	req := ListDocumentsRequest{
		Limit:                  100,
		Offset:                 0,
		IncludeLastViewedPanel: true,
	}
	docs, err := c.fetchDocuments(req)
	if err != nil {
		return nil, err
	}
	return c.findDocByID(docs, id)
}

// GetTranscript fetches the transcript for a document
func (c *Client) GetTranscript(docID string) ([]Utterance, error) {
	payload := map[string]string{"document_id": docID}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/v1/get-document-transcript", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transcript: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result TranscriptResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode transcript: %w", err)
	}

	return result.Utterances, nil
}

// GetDocumentLists fetches all document folders/lists
func (c *Client) GetDocumentLists() ([]DocumentList, error) {
	httpReq, err := http.NewRequest("POST", c.baseURL+"/v2/get-document-lists", bytes.NewReader([]byte("{}")))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch document lists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result DocumentListsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode document lists: %w", err)
	}

	return result.Lists, nil
}

func (c *Client) fetchDocuments(req ListDocumentsRequest) ([]Document, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.baseURL+"/v2/get-documents", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch documents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result ListDocumentsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Docs, nil
}

func (c *Client) findDocByID(docs []Document, id string) (*Document, error) {
	for _, doc := range docs {
		if doc.ID == id || strings.HasPrefix(doc.ID, id) {
			return &doc, nil
		}
	}
	return nil, fmt.Errorf("document not found: %s", id)
}

func (c *Client) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Granola/"+clientVersion)
	req.Header.Set("X-Client-Version", clientVersion)
}
