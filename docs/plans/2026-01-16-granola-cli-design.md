# Granola CLI Design

## Overview

A Go CLI for interacting with Granola meeting notes. Designed for both human use (future TUI) and programmatic access (Claude skills/commands).

Based on reverse-engineering documented at: https://rez0.blog/hacking/2025/05/08/granola-to-obsidian.html

## Goals

1. **CLI-first** - Scriptable commands with JSON output for Claude integration
2. **TUI later** - Interactive Bubble Tea interface (phase 2)
3. **Read operations** - List, show, search notes
4. **Organize later** - Tagging/metadata if API supports it (phase 3)

## Project Structure

```
granola_cli/
├── go.mod
├── go.sum
├── main.go                 # Entry point
├── cmd/
│   ├── root.go             # Root command, global flags
│   ├── list.go             # granola list
│   ├── show.go             # granola show <id>
│   └── search.go           # granola search <query>
├── internal/
│   ├── api/
│   │   ├── client.go       # Granola API client
│   │   └── types.go        # Response structs
│   └── auth/
│       └── credentials.go  # Load tokens from supabase.json
└── tests/
    └── ...
```

## Dependencies

- `github.com/spf13/cobra` - CLI framework

Future (TUI phase):
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `github.com/charmbracelet/bubbles` - Components
- `github.com/charmbracelet/glamour` - Markdown rendering

## Commands

### CLI Mode

```
granola list                     # List recent notes (default: 20)
granola list --limit 50          # List more
granola list --json              # JSON output

granola show <id>                # Show note as markdown
granola show <id> --json         # JSON output
granola show <id> --transcript   # Include transcript if available

granola search <query>           # Search notes by title
granola search <query> --json    # JSON output
```

### Global Flags

- `--credentials <path>` - Override default credentials file
- `--json` - Force JSON output

### Output Formats

**List (text):**
```
ID          Created      Title
─────────────────────────────────────────────────
abc123      2025-05-08   Weekly sync with design team
def456      2025-05-07   Q2 planning session
```

**List (JSON):**
```json
[
  {
    "id": "abc123",
    "title": "Weekly sync with design team",
    "created_at": "2025-05-08T10:30:00Z"
  }
]
```

**Show (text):**
```
# Weekly sync with design team
Created: 2025-05-08 10:30 AM

## Summary
[markdown content]
```

**Show (JSON):**
```json
{
  "id": "abc123",
  "title": "Weekly sync with design team",
  "created_at": "2025-05-08T10:30:00Z",
  "content": "## Summary\n..."
}
```

## Authentication

**Default credentials path:**
```
~/Library/Application Support/Granola/supabase.json
```

**File structure:**
```json
{
  "workos_tokens": "{\"access_token\": \"...\", ...}"
}
```

**Resolution order:**
1. `--credentials` flag
2. `GRANOLA_CREDENTIALS` env var
3. Default path

**Error handling:**
- Clear error if file missing
- Clear error if token not found in file
- No token refresh (user re-opens Granola app if expired)

## API Client

**Base URL:** `https://api.granola.ai`

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
User-Agent: Granola/5.354.0
X-Client-Version: 5.354.0
```

**Endpoints:**

`POST /v2/get-documents`
```json
{
  "limit": 100,
  "offset": 0,
  "include_last_viewed_panel": true
}
```

**Response structure:**
```json
{
  "docs": [
    {
      "id": "...",
      "title": "...",
      "created_at": "...",
      "updated_at": "...",
      "last_viewed_panel": {
        "content": { /* ProseMirror JSON */ }
      }
    }
  ]
}
```

**Search:** Client-side filtering on title (no known search endpoint).

## ProseMirror to Markdown Conversion

The `last_viewed_panel.content` is ProseMirror JSON. Convert to markdown:

- `heading` (level 1-6) → `#` to `######`
- `paragraph` → plain text with newlines
- `bulletList` / `listItem` → `- item`
- `text` → raw text

## Future: TUI Mode (Phase 2)

```
granola                          # Launch interactive TUI
```

**Layout:**
```
┌─ Notes ─────────────────┬─ Preview ─────────────────────┐
│ > Weekly sync with...   │ # Weekly sync with design     │
│   Q2 planning session   │                               │
│   Roadmap review        │ ## Summary                    │
└─────────────────────────┴───────────────────────────────┘
```

**Keybindings (hybrid vim + arrow):**
| Key | Action |
|-----|--------|
| `j` / `↓` | Move down |
| `k` / `↑` | Move up |
| `/` | Start search |
| `Enter` | Full-screen view |
| `Esc` | Back / clear search |
| `q` | Quit |

## Future: Organization Features (Phase 3)

If the Granola API supports updates:
- Add tags/categories to notes
- Store metadata in Granola directly
