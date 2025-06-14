# `md` - Terminal Markdown Viewer

A fast, lightweight CLI utility for rendering Markdown files with rich formatting directly in the terminal. Built in Go with syntax highlighting, and vim-style navigation.

## Features

- **Rich Markdown Rendering**: Support for all standard Markdown elements (headers, lists, tables, links, blockquotes, etc.)
- **Syntax Highlighting**: Code blocks with language-specific highlighting using Chroma
- **Vim Navigation**: Optional vim-style navigation with `less`-like interface
- **Theme Detection**: Automatic terminal theme detection (light/dark)

## Installation

```bash
go install github.com/codehakase/md
```

Or build from source:

```bash
git clone https://github.com/codehakase/md.git
cd md
go build -o md .
./md -v <file.md>
```

## Usage

```
Usage:
  md [flags] <markdown-file>

Flags:
  -h, --help   help for md
  -v, --vim    Enable vim-style navigation
```


### Vim Navigation Keys

When using `--vim` mode, you can navigate using:

- `j` / `k` - Move down/up
- `gg` - Go to top
- `G` - Go to bottom  
- `/` - Search
- `n` - Next search result
- `q` - Quit

## Supported Markdown Features

- **Headers** (`#`, `##`, etc.) with colored styling
- **Text formatting** (bold, italic, strikethrough)
- **Code blocks** with syntax highlighting for 25+ languages
- **Inline code** with theme-appropriate styling
- **Lists** (ordered and unordered) with proper indentation
- **Tables** with borders and header highlighting
- **Links** with URL display
- **Blockquotes** with pipe character styling
- **Task lists** with checkbox rendering
- **Horizontal rules**
