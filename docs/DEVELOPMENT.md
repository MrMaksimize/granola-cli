# Development Notes

## API Discovery

The Granola API was reverse-engineered by proxying the desktop app (documented in [rez0's blog post](https://rez0.blog/hacking/2025/05/08/granola-to-obsidian.html)).

### Authentication

Credentials are stored at:
```
~/Library/Application Support/Granola/supabase.json
```

The file contains a `workos_tokens` field (JSON string) with an `access_token` inside.

### API Endpoint

**Base URL:** `https://api.granola.ai`

**POST /v2/get-documents**

Request:
```json
{
  "limit": 100,
  "offset": 0,
  "include_last_viewed_panel": false
}
```

Headers:
```
Authorization: Bearer <token>
Content-Type: application/json
User-Agent: Granola/5.354.0
X-Client-Version: 5.354.0
```

### Document Structure

Key fields returned per document:

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | UUID |
| `title` | string | Meeting title |
| `created_at` | timestamp | Creation time |
| `updated_at` | timestamp | Last update |
| `notes_markdown` | string | Notes as markdown |
| `notes_plain` | string | Notes as plain text |
| `summary` | string | AI-generated summary |

The `last_viewed_panel.content` field can be either:
- ProseMirror JSON (structured document)
- HTML string (rendered content)

We use `notes_markdown` for simplicity.

## Project Structure

```
granola_cli/
├── main.go                 # Entry point
├── cmd/
│   ├── root.go             # Root command, global flags
│   ├── list.go             # granola list
│   ├── show.go             # granola show <id>
│   └── search.go           # granola search <query>
├── internal/
│   ├── api/
│   │   ├── client.go       # HTTP client
│   │   └── types.go        # API types
│   └── auth/
│       └── credentials.go  # Token loading
└── docs/
    ├── plans/              # Design docs
    └── DEVELOPMENT.md      # This file
```

## Building

```bash
# Build locally
go build -o granola .

# Install to GOPATH/bin
go install .
```

## Commands

```bash
# List recent notes
granola list
granola list --limit 50
granola list --json

# Show a note (supports short ID prefix)
granola show 7cc2c937
granola show 7cc2c937-9746-4514-813e-af2f109d6c71
granola show 7cc2c937 --json

# Search by title
granola search "weekly sync"
granola search "weekly" --json
```

## Global Flags

- `--credentials <path>` - Override default credentials file
- `--json` - Output as JSON

## Future Work

### Phase 2: TUI Mode
- Interactive terminal UI with Bubble Tea
- Split pane: note list + preview
- Vim-style keybindings (j/k, /, Enter, q)

### Phase 3: Organization
- Tagging/categorization if API supports updates
- Local metadata storage as fallback

## Known Limitations

1. **No single-document endpoint** - `show` fetches all docs and filters
2. **Search is client-side** - No known search API, filtering by title only
3. **Token refresh** - Not handled; user must re-open Granola app if expired
4. **macOS only** - Credentials path is hardcoded for macOS
