package renderer

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/codehakase/md/internal/theme"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// CodeHighlighter interface for syntax highlighting
type CodeHighlighter interface {
	Highlight(code, language string) (string, error)
	HighlightInlineCode(code string) string
}

// Renderer renders markdown to styled terminal output
type Renderer struct {
	themeManager *theme.ThemeManager
	goldmark     goldmark.Markdown
}

// New creates a new markdown renderer
func New(themeManager *theme.ThemeManager) *Renderer {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Table,
			extension.Strikethrough,
			extension.TaskList,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	return &Renderer{
		themeManager: themeManager,
		goldmark:     md,
	}
}

// RenderFile renders a markdown file to styled terminal output
func (r *Renderer) RenderFile(filename string, highlighter CodeHighlighter) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return r.RenderContent(content, highlighter)
}

// RenderContent renders markdown content to styled terminal output
func (r *Renderer) RenderContent(content []byte, highlighter CodeHighlighter) (string, error) {
	doc := r.goldmark.Parser().Parse(text.NewReader(content))

	termRenderer := &terminalRenderer{
		themeManager: r.themeManager,
		highlighter:  highlighter,
	}

	var buf bytes.Buffer
	err := termRenderer.render(&buf, content, doc)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	result := buf.String()
	result = TrimTrailingWhitespace(result)
	result = EnsureTrailingNewline(result)

	return result, nil
}

// terminalRenderer handles the actual rendering to terminal format
type terminalRenderer struct {
	themeManager *theme.ThemeManager
	highlighter  CodeHighlighter
}

// render renders the AST node to the writer
func (tr *terminalRenderer) render(w io.Writer, source []byte, node ast.Node) error {
	return ast.Walk(node, func(node ast.Node, entering bool) (ast.WalkStatus, error) {
		err := tr.renderNode(w, source, node, entering)
		if err != nil {
			return ast.WalkStop, err
		}
		return ast.WalkContinue, nil
	})
}

// renderNode renders a specific AST node
func (tr *terminalRenderer) renderNode(w io.Writer, source []byte, node ast.Node, entering bool) error {
	switch n := node.(type) {
	case *ast.Document:
		return tr.renderDocument(w, source, n, entering)
	case *ast.Heading:
		return tr.renderHeading(w, source, n, entering)
	case *ast.Paragraph:
		return tr.renderParagraph(w, source, n, entering)
	case *ast.Text:
		return tr.renderText(w, source, n, entering)
	case *ast.Emphasis:
		return tr.renderEmphasis(w, source, n, entering)
	case *ast.CodeSpan:
		return tr.renderCodeSpan(w, source, n, entering)
	case *ast.FencedCodeBlock:
		return tr.renderFencedCodeBlock(w, source, n, entering)
	case *ast.CodeBlock:
		return tr.renderCodeBlock(w, source, n, entering)
	case *ast.Blockquote:
		return tr.renderBlockquote(w, source, n, entering)
	case *ast.List:
		return tr.renderList(w, source, n, entering)
	case *ast.ListItem:
		return tr.renderListItem(w, source, n, entering)
	case *ast.Link:
		return tr.renderLink(w, source, n, entering)
	case *ast.AutoLink:
		return tr.renderAutoLink(w, source, n, entering)
	case *ast.RawHTML:
		return tr.renderRawHTML(w, source, n, entering)
	case *ast.ThematicBreak:
		return tr.renderThematicBreak(w, source, n, entering)
	case *ast.HTMLBlock:
		return tr.renderHTMLBlock(w, source, n, entering)
	default:
		kind := node.Kind().String()
		switch kind {
		case "Table":
			return tr.renderTable(w, source, node, entering)
		case "TableHeader":
			return tr.renderTableHeader(w, source, node, entering)
		case "TableRow":
			return tr.renderTableRow(w, source, node, entering)
		case "TableCell":
			return tr.renderTableCell(w, source, node, entering)
		case "Strikethrough":
			return tr.renderStrikethrough(w, source, node, entering)
		case "TaskCheckBox":
			return tr.renderTaskCheckBox(w, source, node, entering)
		}
		return nil
	}
}

func (tr *terminalRenderer) renderDocument(w io.Writer, source []byte, n *ast.Document, entering bool) error {
	return nil
}

func (tr *terminalRenderer) renderHeading(w io.Writer, source []byte, n *ast.Heading, entering bool) error {
	if entering {
		var headerTheme theme.ColorKey
		var prefix string

		switch n.Level {
		case 1:
			headerTheme = theme.Header1
			prefix = H1Prefix
		case 2:
			headerTheme = theme.Header2
			prefix = H2Prefix
		case 3:
			headerTheme = theme.Header3
			prefix = H3Prefix
		case 4:
			headerTheme = theme.Header4
			prefix = H4Prefix
		case 5:
			headerTheme = theme.Header5
			prefix = H5Prefix
		default:
			headerTheme = theme.Header6
			prefix = H6Prefix
		}

		if n.PreviousSibling() != nil {
			fmt.Fprint(w, "\n")
		}

		fmt.Fprint(w, tr.themeManager.StyleNoReset(prefix, headerTheme))
	} else {
		fmt.Fprint(w, tr.themeManager.Reset())
		fmt.Fprint(w, "\n")
	}
	return nil
}

func (tr *terminalRenderer) renderParagraph(w io.Writer, source []byte, n *ast.Paragraph, entering bool) error {
	if !entering {
		fmt.Fprint(w, "\n")
		if n.NextSibling() != nil {
			fmt.Fprint(w, "\n")
		}
	}
	return nil
}

func (tr *terminalRenderer) renderText(w io.Writer, source []byte, n *ast.Text, entering bool) error {
	if entering {
		// Skip text rendering if we're inside a CodeSpan (already handled by renderCodeSpan)
		if _, isCodeSpanChild := n.Parent().(*ast.CodeSpan); isCodeSpanChild {
			return nil
		}

		segment := n.Segment
		value := segment.Value(source)

		if n.IsRaw() {
			fmt.Fprint(w, string(value))
		} else {
			lines := strings.Split(string(value), "\n")
			for i, line := range lines {
				if i > 0 {
					if n.HardLineBreak() || (i == len(lines)-1 && n.SoftLineBreak()) {
						fmt.Fprint(w, "\n")
					} else {
						fmt.Fprint(w, " ")
					}
				}
				fmt.Fprint(w, line)
			}
		}
	}
	return nil
}

func (tr *terminalRenderer) renderEmphasis(w io.Writer, source []byte, n *ast.Emphasis, entering bool) error {
	if entering {
		if n.Level == 2 {
			fmt.Fprint(w, tr.themeManager.GetColor(theme.Bold))
		} else {
			fmt.Fprint(w, tr.themeManager.GetColor(theme.Italic))
		}
	} else {
		fmt.Fprint(w, tr.themeManager.Reset())
	}
	return nil
}

func (tr *terminalRenderer) renderCodeSpan(w io.Writer, source []byte, n *ast.CodeSpan, entering bool) error {
	if entering {
		value := string(n.Text(source))
		highlighted := tr.highlighter.HighlightInlineCode(value)
		fmt.Fprint(w, highlighted)
	}
	return nil
}

func (tr *terminalRenderer) renderFencedCodeBlock(w io.Writer, source []byte, n *ast.FencedCodeBlock, entering bool) error {
	if entering {
		language := ""
		if n.Info != nil {
			language = string(n.Info.Text(source))
		}

		var code strings.Builder
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			code.Write(line.Value(source))
		}

		highlighted, err := tr.highlighter.Highlight(code.String(), language)
		if err != nil {
			highlighted = tr.themeManager.Style(code.String(), theme.Code)
		}

		fmt.Fprint(w, "\n")
		fmt.Fprint(w, Indent(strings.TrimRight(highlighted, "\n"), 1))
		fmt.Fprint(w, "\n\n")
	}
	return nil
}

func (tr *terminalRenderer) renderCodeBlock(w io.Writer, source []byte, n *ast.CodeBlock, entering bool) error {
	if entering {
		var code strings.Builder
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			code.Write(line.Value(source))
		}

		styled := tr.themeManager.Style(code.String(), theme.Code)

		fmt.Fprint(w, "\n")
		fmt.Fprint(w, Indent(strings.TrimRight(styled, "\n"), 1))
		fmt.Fprint(w, "\n\n")
	}
	return nil
}

func (tr *terminalRenderer) renderBlockquote(w io.Writer, source []byte, n *ast.Blockquote, entering bool) error {
	if entering {
		fmt.Fprint(w, "\n")
		fmt.Fprint(w, tr.themeManager.GetColor(theme.BlockQuote))
		fmt.Fprint(w, "│ ")
	} else {
		fmt.Fprint(w, tr.themeManager.Reset())
		fmt.Fprint(w, "\n\n")
	}
	return nil
}

func (tr *terminalRenderer) renderList(w io.Writer, source []byte, n *ast.List, entering bool) error {
	if entering {
		if n.PreviousSibling() != nil {
			fmt.Fprint(w, "\n")
		}
	} else {
		if n.NextSibling() != nil {
			fmt.Fprint(w, "\n")
		}
	}
	return nil
}

func (tr *terminalRenderer) renderListItem(w io.Writer, source []byte, n *ast.ListItem, entering bool) error {
	if entering {
		level := 0
		parent := n.Parent()
		for parent != nil {
			if _, isList := parent.(*ast.List); isList {
				level++
			}
			parent = parent.Parent()
		}
		level--

		indent := strings.Repeat("  ", level)

		var marker string
		if list, ok := n.Parent().(*ast.List); ok {
			if list.IsOrdered() {
				// Using generic numbered marker for simplicity - would need counter for proper numbering
				marker = tr.themeManager.Style("1.", theme.OrderedList)
			} else {
				marker = tr.themeManager.Style("•", theme.BulletPoint)
			}
		} else {
			marker = tr.themeManager.Style("•", theme.BulletPoint)
		}

		fmt.Fprintf(w, "\u00A0%s%s ", indent, marker)
	} else {
		fmt.Fprint(w, "\n")
	}
	return nil
}

func (tr *terminalRenderer) renderLink(w io.Writer, source []byte, n *ast.Link, entering bool) error {
	if entering {
		fmt.Fprint(w, tr.themeManager.GetColor(theme.Link))
	} else {
		url := string(n.Destination)
		fmt.Fprintf(w, " (%s)", url)
		fmt.Fprint(w, tr.themeManager.Reset())
	}
	return nil
}

func (tr *terminalRenderer) renderAutoLink(w io.Writer, source []byte, n *ast.AutoLink, entering bool) error {
	if entering {
		url := string(n.URL(source))
		fmt.Fprint(w, tr.themeManager.Style(url, theme.Link))
	}
	return nil
}

func (tr *terminalRenderer) renderRawHTML(w io.Writer, source []byte, n *ast.RawHTML, entering bool) error {
	return nil
}

func (tr *terminalRenderer) renderHTMLBlock(w io.Writer, source []byte, n *ast.HTMLBlock, entering bool) error {
	return nil
}

func (tr *terminalRenderer) renderThematicBreak(w io.Writer, source []byte, n *ast.ThematicBreak, entering bool) error {
	if entering {
		fmt.Fprint(w, "\n")
		rule := strings.Repeat("─", 50)
		fmt.Fprint(w, tr.themeManager.Style(rule, theme.TableBorder))
		fmt.Fprint(w, "\n\n")
	}
	return nil
}

// Extension-specific renderers
func (tr *terminalRenderer) renderTable(w io.Writer, source []byte, node ast.Node, entering bool) error {
	if entering {
		fmt.Fprint(w, "\n")
	} else {
		fmt.Fprint(w, "\n")
	}
	return nil
}

func (tr *terminalRenderer) renderTableHeader(w io.Writer, source []byte, node ast.Node, entering bool) error {
	return nil
}

func (tr *terminalRenderer) renderTableRow(w io.Writer, source []byte, node ast.Node, entering bool) error {
	if entering {
		fmt.Fprint(w, tr.themeManager.GetColor(theme.TableBorder))
		fmt.Fprint(w, "│")
	} else {
		fmt.Fprint(w, tr.themeManager.Reset())
		fmt.Fprint(w, "\n")

		if node.Kind().String() == "TableHeader" {
			cellCount := 0
			for child := node.FirstChild(); child != nil; child = child.NextSibling() {
				cellCount++
			}

			fmt.Fprint(w, tr.themeManager.GetColor(theme.TableBorder))
			for i := 0; i < cellCount; i++ {
				fmt.Fprint(w, "├")
				fmt.Fprint(w, strings.Repeat("─", 15))
				if i < cellCount-1 {
					fmt.Fprint(w, "┼")
				}
			}
			fmt.Fprint(w, "┤")
			fmt.Fprint(w, tr.themeManager.Reset())
			fmt.Fprint(w, "\n")
		}
	}
	return nil
}

func (tr *terminalRenderer) renderTableCell(w io.Writer, source []byte, node ast.Node, entering bool) error {
	if entering {
		fmt.Fprint(w, " ")
		if node.Parent().Kind().String() == "TableHeader" {
			fmt.Fprint(w, tr.themeManager.GetColor(theme.TableHeader))
		}
	} else {
		if node.Parent().Kind().String() == "TableHeader" {
			fmt.Fprint(w, tr.themeManager.Reset())
		}

		fmt.Fprint(w, strings.Repeat(" ", 14))
		fmt.Fprint(w, tr.themeManager.GetColor(theme.TableBorder))
		fmt.Fprint(w, "│")
	}
	return nil
}

func (tr *terminalRenderer) renderStrikethrough(w io.Writer, source []byte, node ast.Node, entering bool) error {
	if entering {
		fmt.Fprint(w, tr.themeManager.GetColor(theme.Strikethrough))
	} else {
		fmt.Fprint(w, tr.themeManager.Reset())
	}
	return nil
}

func (tr *terminalRenderer) renderTaskCheckBox(w io.Writer, source []byte, node ast.Node, entering bool) error {
	if entering {
		// Simplified implementation - would need to cast to extension type to check if checked
		fmt.Fprint(w, tr.themeManager.Style("[ ]", theme.BulletPoint))
		fmt.Fprint(w, " ")
	}
	return nil
}

