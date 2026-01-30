package format

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ProseMirrorNode represents a node in the ProseMirror document tree
type ProseMirrorNode struct {
	Type    string            `json:"type"`
	Content []ProseMirrorNode `json:"content,omitempty"`
	Text    string            `json:"text,omitempty"`
	Marks   []Mark            `json:"marks,omitempty"`
	Attrs   map[string]any    `json:"attrs,omitempty"`
}

// Mark represents a text mark (bold, italic, link, etc.)
type Mark struct {
	Type  string         `json:"type"`
	Attrs map[string]any `json:"attrs,omitempty"`
}

// ProseMirrorToMarkdown converts ProseMirror JSON to markdown
func ProseMirrorToMarkdown(data json.RawMessage) (string, error) {
	var doc ProseMirrorNode
	if err := json.Unmarshal(data, &doc); err != nil {
		return "", fmt.Errorf("failed to parse ProseMirror JSON: %w", err)
	}

	var sb strings.Builder
	renderNodes(&sb, doc.Content, 0)
	return strings.TrimRight(sb.String(), "\n") + "\n", nil
}

func renderNodes(sb *strings.Builder, nodes []ProseMirrorNode, depth int) {
	for i, node := range nodes {
		renderNode(sb, node, depth, i)
	}
}

func renderNode(sb *strings.Builder, node ProseMirrorNode, depth int, index int) {
	switch node.Type {
	case "doc":
		renderNodes(sb, node.Content, depth)

	case "paragraph":
		renderInline(sb, node.Content)
		sb.WriteString("\n\n")

	case "heading":
		level := 1
		if l, ok := node.Attrs["level"]; ok {
			if lf, ok := l.(float64); ok {
				level = int(lf)
			}
		}
		sb.WriteString(strings.Repeat("#", level) + " ")
		renderInline(sb, node.Content)
		sb.WriteString("\n\n")

	case "bulletList":
		renderListItems(sb, node.Content, depth, false)

	case "orderedList":
		renderListItems(sb, node.Content, depth, true)

	case "listItem":
		renderListItem(sb, node, depth, false, index)

	case "blockquote":
		var inner strings.Builder
		renderNodes(&inner, node.Content, depth)
		for _, line := range strings.Split(strings.TrimRight(inner.String(), "\n"), "\n") {
			sb.WriteString("> " + line + "\n")
		}
		sb.WriteString("\n")

	case "codeBlock":
		lang := ""
		if l, ok := node.Attrs["language"]; ok {
			if ls, ok := l.(string); ok {
				lang = ls
			}
		}
		sb.WriteString("```" + lang + "\n")
		renderInline(sb, node.Content)
		sb.WriteString("\n```\n\n")

	case "hardBreak":
		sb.WriteString("\n")

	case "horizontalRule":
		sb.WriteString("---\n\n")

	case "text":
		sb.WriteString(renderMarkedText(node))

	default:
		// Unknown node types: render children if any
		renderNodes(sb, node.Content, depth)
	}
}

func renderListItems(sb *strings.Builder, items []ProseMirrorNode, depth int, ordered bool) {
	for i, item := range items {
		renderListItem(sb, item, depth, ordered, i)
	}
	if depth == 0 {
		sb.WriteString("\n")
	}
}

func renderListItem(sb *strings.Builder, item ProseMirrorNode, depth int, ordered bool, index int) {
	indent := strings.Repeat("  ", depth)
	prefix := "- "
	if ordered {
		prefix = fmt.Sprintf("%d. ", index+1)
	}

	for j, child := range item.Content {
		switch child.Type {
		case "paragraph":
			if j == 0 {
				sb.WriteString(indent + prefix)
			} else {
				sb.WriteString(indent + strings.Repeat(" ", len(prefix)))
			}
			renderInline(sb, child.Content)
			sb.WriteString("\n")
		case "bulletList":
			renderListItems(sb, child.Content, depth+1, false)
		case "orderedList":
			renderListItems(sb, child.Content, depth+1, true)
		default:
			if j == 0 {
				sb.WriteString(indent + prefix)
			}
			renderNode(sb, child, depth, j)
		}
	}
}

func renderInline(sb *strings.Builder, nodes []ProseMirrorNode) {
	for _, node := range nodes {
		switch node.Type {
		case "text":
			sb.WriteString(renderMarkedText(node))
		case "hardBreak":
			sb.WriteString("\n")
		default:
			renderNode(sb, node, 0, 0)
		}
	}
}

func renderMarkedText(node ProseMirrorNode) string {
	text := node.Text
	for _, mark := range node.Marks {
		switch mark.Type {
		case "bold", "strong":
			text = "**" + text + "**"
		case "italic", "em":
			text = "_" + text + "_"
		case "code":
			text = "`" + text + "`"
		case "link":
			href := ""
			if h, ok := mark.Attrs["href"]; ok {
				if hs, ok := h.(string); ok {
					href = hs
				}
			}
			text = "[" + text + "](" + href + ")"
		case "strikethrough":
			text = "~~" + text + "~~"
		}
	}
	return text
}
