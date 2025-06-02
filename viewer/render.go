package viewer

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// updateViewNodes collects all visible nodes based on expansion state and filters
func (m *Model) updateViewNodes() {
	m.viewNodes = nil
	m.collectViewNodes(m.root)

	if m.filter != "" && !m.jsonpathMode {
		filtered := make([]*Node, 0)
		filterLower := strings.ToLower(m.filter)
		for _, node := range m.viewNodes {
			if strings.Contains(strings.ToLower(node.Path), filterLower) ||
				strings.Contains(strings.ToLower(node.Key), filterLower) ||
				(node.Value != nil && strings.Contains(strings.ToLower(fmt.Sprintf("%v", node.Value)), filterLower)) {
				filtered = append(filtered, node)
			}
		}
		m.viewNodes = filtered
	}

	if m.cursor >= len(m.viewNodes) && len(m.viewNodes) > 0 {
		m.cursor = len(m.viewNodes) - 1
	}
	if m.cursor < 0 {
		m.cursor = 0
	}
}

// collectViewNodes recursively collects visible nodes
func (m *Model) collectViewNodes(node *Node) {
	m.viewNodes = append(m.viewNodes, node)
	if node.Expanded {
		for _, child := range node.Children {
			m.collectViewNodes(child)
		}
	}
}

// updateViewport updates the viewport content and positioning
func (m *Model) updateViewport() {
	if len(m.viewNodes) == 0 {
		m.viewport.SetContent("")
		return
	}

	var content strings.Builder
	for i, node := range m.viewNodes {
		line := m.renderNode(node, i == m.cursor)
		content.WriteString(line + "\n")
	}

	m.viewport.SetContent(content.String())

	if m.cursor >= 0 && m.cursor < len(m.viewNodes) {
		m.viewport.SetYOffset(max(0, m.cursor-m.viewport.Height/2))
	}
}

// renderNode renders a single tree node
func (m Model) renderNode(node *Node, isCursor bool) string {
	depth := node.GetDepth()
	indent := strings.Repeat("  ", depth)

	var icon string
	if len(node.Children) > 0 {
		if node.Expanded {
			icon = "▼ "
		} else {
			icon = "▶ "
		}
	} else {
		icon = "  "
	}

	var keyPart string
	if node.Key != "" {
		keyPart = m.config.Theme.Key.Render(node.Key) + ": "
	}

	var valuePart string
	switch node.Type {
	case ObjectNode:
		if node.Expanded {
			valuePart = "{"
		} else {
			valuePart = fmt.Sprintf("{...} (%d items)", len(node.Children))
		}
	case ArrayNode:
		if node.Expanded {
			valuePart = "["
		} else {
			valuePart = fmt.Sprintf("[...] (%d items)", len(node.Children))
		}
	case StringNode:
		valuePart = m.config.Theme.String.Render(fmt.Sprintf("\"%s\"", node.Value))
	case NumberNode:
		valuePart = m.config.Theme.Number.Render(fmt.Sprintf("%v", node.Value))
	case BoolNode:
		valuePart = m.config.Theme.Bool.Render(fmt.Sprintf("%v", node.Value))
	case NullNode:
		valuePart = m.config.Theme.Null.Render("null")
	}

	// Check if this node is a search match
	isMatch := false
	if len(m.searchMatches) > 0 {
		for _, match := range m.searchMatches {
			if match == node {
				isMatch = true
				break
			}
		}
	}

	line := indent + icon + keyPart + valuePart

	if isMatch {
		line = m.config.Theme.Match.Render(line)
	} else if isCursor {
		line = m.config.Theme.Cursor.Render(line)
	}

	return line
}

// renderHeader renders the header section (only in non-embedded mode)
func (m Model) renderHeader() string {
	if m.embedded {
		return ""
	}

	title := m.config.Theme.Header.Render("Bonsai")
	if m.filename != "" {
		title = m.config.Theme.Header.Render(fmt.Sprintf("Bonsai - %s", m.filename))
	}

	stats := ""
	if m.fileSize > 0 {
		stats = m.config.Theme.Status.Render(fmt.Sprintf("Size: %.1fKB | Nodes: %d",
			float64(m.fileSize)/1024, m.nodeCount))
	} else {
		stats = m.config.Theme.Status.Render(fmt.Sprintf("Nodes: %d", m.nodeCount))
	}

	var filterInfo string
	if m.filterMode {
		filterInfo = m.config.Theme.Filter.Render(fmt.Sprintf("Filter: %s█", m.filter))
	} else if m.jsonpathMode {
		filterInfo = m.config.Theme.JSONPath.Render(fmt.Sprintf("JSONPath: %s█", m.filter))
	} else if m.searchMode {
		searchInfo := fmt.Sprintf("Search: %s█", m.filter)
		if len(m.searchMatches) > 0 {
			searchInfo = fmt.Sprintf("Search: %s█ (%d/%d)", m.filter, m.searchIndex+1, len(m.searchMatches))
		}
		filterInfo = m.config.Theme.Search.Render(searchInfo)
	} else if m.gotoMode {
		filterInfo = m.config.Theme.Goto.Render(fmt.Sprintf("Goto: %s█", m.filter))
	} else if m.filter != "" {
		if strings.HasPrefix(m.filter, "$") {
			filterInfo = m.config.Theme.JSONPath.Render(fmt.Sprintf("Active JSONPath: %s", m.filter))
		} else {
			filterInfo = m.config.Theme.Filter.Render(fmt.Sprintf("Active Filter: %s", m.filter))
		}
	}

	breadcrumb := ""
	if m.cursor < len(m.viewNodes) && len(m.viewNodes) > 0 {
		node := m.viewNodes[m.cursor]
		breadcrumb = m.config.Theme.Breadcrumb.Render(fmt.Sprintf("Path: %s", node.Path))
	}

	headerLine1 := lipgloss.JoinHorizontal(lipgloss.Left, title, strings.Repeat(" ", 5), stats)
	headerLine2 := ""
	if filterInfo != "" {
		headerLine2 = filterInfo
	} else {
		headerLine2 = breadcrumb
	}

	return lipgloss.JoinVertical(lipgloss.Left, headerLine1, headerLine2, "")
}

// renderFooter renders the footer section (only in non-embedded mode)
func (m Model) renderFooter() string {
	if m.embedded {
		return ""
	}

	if m.showHelp && m.config.ShowHelp {
		return m.renderManualHelp()
	}

	if m.filterMode {
		return "\n" + m.config.Theme.Filter.Render("Press Enter to apply filter, Esc to cancel")
	} else if m.jsonpathMode {
		return "\n" + m.config.Theme.JSONPath.Render("Press Enter to apply JSONPath, Esc to cancel")
	} else if m.searchMode {
		return "\n" + m.config.Theme.Search.Render("Press Enter to search, Esc to cancel")
	} else if m.gotoMode {
		return "\n" + m.config.Theme.Goto.Render("Press Enter to goto path, Esc to cancel")
	}

	helpText := m.config.Theme.Status.Render("Press ? for help")
	return "\n" + helpText
}

// renderManualHelp renders the help overlay
func (m Model) renderManualHelp() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)

	help := strings.Builder{}
	help.WriteString("\n")
	help.WriteString(titleStyle.Render("━━━ Bonsai Help ━━━") + "\n\n")

	help.WriteString(helpStyle.Render("Navigation:") + "\n")
	help.WriteString("  ↑/k, ↓/j, ←/h, →/l     Navigate tree\n")
	help.WriteString("  PgUp/Ctrl+U, PgDn/Ctrl+D   Page up/down\n")
	help.WriteString("  Home/g, End/G           Go to top/bottom\n\n")

	help.WriteString(helpStyle.Render("Tree Operations:") + "\n")
	help.WriteString("  Enter/Space/l           Expand/collapse node\n")
	help.WriteString("  h                       Collapse or go to parent\n")
	help.WriteString("  E, C                    Expand/collapse all\n\n")

	help.WriteString(helpStyle.Render("Search & Filter:") + "\n")
	help.WriteString("  /                       Text filter\n")
	help.WriteString("  $                       JSONPath query\n")
	help.WriteString("  s/Ctrl+F                Search content\n")
	help.WriteString("  :/Ctrl+G                Goto path\n")
	help.WriteString("  n, N                    Next/prev match\n\n")

	if m.config.EnableClipboard {
		help.WriteString(helpStyle.Render("Clipboard:") + "\n")
		help.WriteString("  c, p, y                 Copy value/path/key\n\n")
	}

	help.WriteString(helpStyle.Render("Utility:") + "\n")
	help.WriteString("  r/Ctrl+R                Reset view\n")
	help.WriteString("  ?, q/Esc                Help, Quit\n\n")

	help.WriteString(titleStyle.Render("Press ? to close help"))

	return help.String()
}