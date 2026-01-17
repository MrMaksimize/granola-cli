# Granola CLI

A command-line interface for [Granola](https://granola.ai) meeting notes.

## Installation

### Prerequisites

- Go 1.21+ ([install](https://go.dev/doc/install))
- Granola desktop app (for authentication)
- `~/go/bin` in your PATH

### Quick Install

```bash
# Install directly from GitHub
go install github.com/MrMaksimize/granola-cli@latest

# The binary is named granola-cli, create alias for convenience
ln -sf ~/go/bin/granola-cli ~/go/bin/granola
```

### From Source

```bash
git clone https://github.com/MrMaksimize/granola-cli.git
cd granola-cli
go install .
ln -sf ~/go/bin/granola-cli ~/go/bin/granola
```

### PATH Setup

Ensure `~/go/bin` is in your PATH. Add to your shell config (`~/.zshrc` or `~/.bashrc`):

```bash
export PATH="$HOME/go/bin:$PATH"
```

Then reload: `source ~/.zshrc`

### Verify Installation

```bash
granola list
```

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
