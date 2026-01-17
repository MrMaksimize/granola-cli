# Granola CLI

A command-line interface for [Granola](https://granola.ai) meeting notes.

## Installation

```bash
# Clone and install
git clone https://github.com/mrm/granola-cli.git
cd granola-cli
go install .
```

Requires Go 1.21+ and the Granola desktop app (for authentication).

## Usage

```bash
# List recent notes
granola list
granola list --limit 50

# Show a note (supports short ID prefix)
granola show 7cc2c937

# Search by title
granola search "weekly sync"

# JSON output (for scripts/Claude)
granola list --json
granola show 7cc2c937 --json
granola search "weekly" --json
```

## Authentication

The CLI uses credentials from the Granola desktop app:

```
~/Library/Application Support/Granola/supabase.json
```

Override with `--credentials <path>` or `GRANOLA_CREDENTIALS` env var.

## Commands

| Command | Description |
|---------|-------------|
| `list` | List recent meeting notes |
| `show <id>` | Show full note content |
| `search <query>` | Search notes by title |

### Global Flags

| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |
| `--credentials <path>` | Custom credentials file |

## Roadmap

- [x] CLI with list/show/search
- [ ] Interactive TUI (Bubble Tea)
- [ ] Organization features (tags, categories)

## Credits

API reverse-engineering based on [rez0's blog post](https://rez0.blog/hacking/2025/05/08/granola-to-obsidian.html).
