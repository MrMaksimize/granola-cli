# Developing Granola CLI

## Prerequisites

- Go 1.21+
- Granola desktop app (for credentials)

## Setup

```bash
git clone https://github.com/mrm/granola-cli.git
cd granola-cli
go mod tidy
```

## Building

```bash
# Build binary
go build -o granola .

# Run directly
go run . list

# Install to GOPATH/bin
go install .
```

## Testing

```bash
# Test list command
./granola list

# Test with JSON
./granola list --json | jq .

# Test show (use an ID from list)
./granola show <id>

# Test search
./granola search "meeting"
```

## Project Structure

```
.
├── main.go                 # Entry point - just calls cmd.Execute()
├── cmd/
│   ├── root.go             # Root command, global flags (--json, --credentials)
│   ├── list.go             # List command with --limit
│   ├── show.go             # Show command (accepts ID or prefix)
│   └── search.go           # Search command (client-side title filter)
├── internal/
│   ├── api/
│   │   ├── client.go       # HTTP client for Granola API
│   │   └── types.go        # Request/response types
│   └── auth/
│       └── credentials.go  # Token loading from supabase.json
└── docs/
    └── plans/              # Design documents
```

## Adding a New Command

1. Create `cmd/<command>.go`:

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycommand",
    Short: "Description",
    RunE:  runMyCommand,
}

func init() {
    rootCmd.AddCommand(myCmd)
}

func runMyCommand(cmd *cobra.Command, args []string) error {
    // Use credentialsPath and jsonOutput from root.go
    token, err := auth.LoadCredentials(credentialsPath)
    if err != nil {
        return err
    }

    client := api.NewClient(token)
    // ... do something

    if jsonOutput {
        // JSON output
    }
    // Text output
    return nil
}
```

2. Build and test: `go build -o granola . && ./granola mycommand`

## API Client

The client is in `internal/api/client.go`. To add a new API method:

```go
func (c *Client) MyMethod() (Result, error) {
    // Use c.httpClient and c.setHeaders(req)
}
```

## Output Patterns

All commands should support dual output:

```go
if jsonOutput {
    enc := json.NewEncoder(os.Stdout)
    enc.SetIndent("", "  ")
    return enc.Encode(data)
}

// Human-readable output
w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
// ...
return w.Flush()
```

## Debugging API

To see raw API responses:

```go
// Temporarily in client.go
respBody, _ := io.ReadAll(resp.Body)
fmt.Println(string(respBody))
```

Or use a proxy like Caido/Burp with `HTTP_PROXY` env var.

## Common Issues

### "credentials file not found"
- Ensure Granola desktop app is installed and you're logged in
- Check `~/Library/Application Support/Granola/supabase.json` exists

### "API returned status 401"
- Token expired; re-open Granola desktop app to refresh

### Empty note content
- Some notes have empty `notes_markdown`; this is a Granola issue, not CLI

## Phase 2: TUI Development

When adding TUI mode:

```bash
# Additional dependencies
go get github.com/charmbracelet/bubbletea
go get github.com/charmbracelet/lipgloss
go get github.com/charmbracelet/bubbles
go get github.com/charmbracelet/glamour
```

TUI code should go in `internal/tui/`:
- `app.go` - Main Bubble Tea model
- `list.go` - Note list component
- `preview.go` - Note preview pane
- `styles.go` - Lip Gloss styles

Launch TUI when `granola` is run with no arguments.
