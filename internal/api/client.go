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
	baseURL       = "https://api.granola.ai"
	clientVersion = "5.354.0"
)

// Client is a Granola API client
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new Granola API client
func NewClient(token string) *Client {
	return &Client{
		token:      token,
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

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", baseURL+"/v2/get-documents", bytes.NewReader(body))
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

// GetDocument fetches a single document by ID
func (c *Client) GetDocument(id string) (*Document, error) {
	// The API doesn't have a single-document endpoint that we know of,
	// so we fetch all and filter. This is inefficient but works.
	// Also support prefix matching for short IDs
	docs, err := c.ListDocuments(100, 0)
	if err != nil {
		return nil, err
	}

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
