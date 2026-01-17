# Claude Code Instructions

## Project Overview

Go CLI for interacting with Granola meeting notes. Designed for both human use (future TUI) and programmatic access (Claude skills/commands).

## Design Philosophy

### Dual-Mode Output
Every command supports both human-readable and JSON output:
- Default: Formatted tables/text for terminal use
- `--json`: Structured JSON for scripts and Claude integration

### Short ID Support
Document IDs are UUIDs but commands accept prefixes:
```bash
granola show 7cc2c937  # Works with first 8 chars
```

### Fail Fast
- Clear error messages when credentials missing
- No silent failures or empty outputs without explanation

### YAGNI
- No token refresh (user re-opens Granola app)
- No caching (API is fast enough)
- No pagination UI (just `--limit` flag)

## Code Standards

### Project Layout
```
cmd/           # Cobra commands (one file per command)
internal/api/  # API client and types
internal/auth/ # Credential handling
```

### Naming
- Files: lowercase, underscore separation
- Types: PascalCase
- Functions: camelCase for private, PascalCase for exported

### Error Handling
- Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- Commands return errors; root command handles exit codes

### Dependencies
- `cobra` for CLI framework
- Standard library for HTTP, JSON
- No external logging framework (just fmt)

## API Notes

### Endpoint
```
POST https://api.granola.ai/v2/get-documents
```

### Key Response Fields
- `notes_markdown` - Use this for content (pre-formatted)
- `notes_plain` - Plain text alternative
- `last_viewed_panel.content` - Can be ProseMirror JSON or HTML string (avoid)

### Authentication
Bearer token from `~/Library/Application Support/Granola/supabase.json`

## Future Work

### Phase 2: TUI
- Bubble Tea for interactive mode
- `granola` (no args) launches TUI
- Commands still work for scripting

### Phase 3: Organization
- Tags/categories if API supports updates
- Local metadata as fallback
