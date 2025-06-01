package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type NodeType int

const (
	ObjectNode NodeType = iota
	ArrayNode
	StringNode
	NumberNode
	BoolNode
	NullNode
)

type Node struct {
	Key      string
	Value    interface{}
	Type     NodeType
	Children []*Node
	Parent   *Node
	Expanded bool
	Path     string
}

type keyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Enter        key.Binding
	PageUp       key.Binding
	PageDown     key.Binding
	Home         key.Binding
	End          key.Binding
	Filter       key.Binding
	JSONPath     key.Binding
	Copy         key.Binding
	CopyPath     key.Binding
	CopyKey      key.Binding
	Reset        key.Binding
	ExpandAll    key.Binding
	CollapseAll  key.Binding
	Goto         key.Binding
	Search       key.Binding
	NextMatch    key.Binding
	PrevMatch    key.Binding
	Help         key.Binding
	Quit         key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right},
		{k.PageUp, k.PageDown, k.Home, k.End},
		{k.Enter, k.ExpandAll, k.CollapseAll},
		{k.Filter, k.JSONPath, k.Search, k.Goto},
		{k.Copy, k.CopyPath, k.CopyKey},
		{k.NextMatch, k.PrevMatch, k.Reset},
		{k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "collapse/parent"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "expand"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "expand/collapse"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("pgup", "ctrl+u"),
		key.WithHelp("pgup/ctrl+u", "page up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("pgdown", "ctrl+d"),
		key.WithHelp("pgdn/ctrl+d", "page down"),
	),
	Home: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("home/g", "go to top"),
	),
	End: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("end/G", "go to bottom"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter text"),
	),
	JSONPath: key.NewBinding(
		key.WithKeys("$"),
		key.WithHelp("$", "jsonpath query"),
	),
	Copy: key.NewBinding(
		key.WithKeys("c"),
		key.WithHelp("c", "copy value"),
	),
	CopyPath: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "copy path"),
	),
	CopyKey: key.NewBinding(
		key.WithKeys("y"),
		key.WithHelp("y", "copy key"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r", "ctrl+r"),
		key.WithHelp("r/ctrl+r", "reset view"),
	),
	ExpandAll: key.NewBinding(
		key.WithKeys("E"),
		key.WithHelp("E", "expand all"),
	),
	CollapseAll: key.NewBinding(
		key.WithKeys("C"),
		key.WithHelp("C", "collapse all"),
	),
	Goto: key.NewBinding(
		key.WithKeys(":", "ctrl+g"),
		key.WithHelp(":/ctrl+g", "goto path"),
	),
	Search: key.NewBinding(
		key.WithKeys("s", "ctrl+f"),
		key.WithHelp("s/ctrl+f", "search"),
	),
	NextMatch: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next match"),
	),
	PrevMatch: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "prev match"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c", "esc"),
		key.WithHelp("q/esc", "quit"),
	),
}

type Model struct {
	root          *Node
	cursor        int
	viewNodes     []*Node
	viewport      viewport.Model
	help          help.Model
	progress      progress.Model
	keys          keyMap
	filename      string
	fileSize      int64
	nodeCount     int
	filter        string
	filterMode    bool
	jsonpathMode  bool
	searchMode    bool
	gotoMode      bool
	showHelp      bool
	rawData       interface{}
	searchMatches []*Node
	searchIndex   int
}

var (
	headerStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	statusStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	keyStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	stringStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("220"))
	numberStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("208"))
	boolStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	nullStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle     = lipgloss.NewStyle().Background(lipgloss.Color("240"))
	filterStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	jsonpathStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("171")).Bold(true)
	searchStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	gotoStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
	breadcrumbStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	matchStyle      = lipgloss.NewStyle().Background(lipgloss.Color("220")).Foreground(lipgloss.Color("16"))
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: json-viewer <file.json>")
		os.Exit(1)
	}

	filename := os.Args[1]
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		log.Fatal(err)
	}

	root := buildTree(jsonData, "", "$")
	root.Expanded = true

	vp := viewport.New(80, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	prog := progress.New(progress.WithDefaultGradient())

	helpModel := help.New()
	helpModel.Width = 80

	m := Model{
		root:      root,
		cursor:    0,
		viewport:  vp,
		help:      helpModel,
		progress:  prog,
		keys:      keys,
		filename:  filepath.Base(filename),
		fileSize:  fileInfo.Size(),
		nodeCount: countNodes(root),
		rawData:   jsonData,
	}
	m.updateViewNodes()
	m.updateViewport()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func buildTree(data interface{}, key, path string) *Node {
	node := &Node{
		Key:  key,
		Path: path,
	}

	switch v := data.(type) {
	case map[string]interface{}:
		node.Type = ObjectNode
		node.Value = v
		for k, val := range v {
			child := buildTree(val, k, path+"."+k)
			child.Parent = node
			node.Children = append(node.Children, child)
		}
	case []interface{}:
		node.Type = ArrayNode
		node.Value = v
		for i, val := range v {
			child := buildTree(val, fmt.Sprintf("[%d]", i), fmt.Sprintf("%s[%d]", path, i))
			child.Parent = node
			node.Children = append(node.Children, child)
		}
	case string:
		node.Type = StringNode
		node.Value = v
	case float64:
		node.Type = NumberNode
		node.Value = v
	case bool:
		node.Type = BoolNode
		node.Value = v
	case nil:
		node.Type = NullNode
		node.Value = nil
	}

	return node
}

func countNodes(node *Node) int {
	count := 1
	for _, child := range node.Children {
		count += countNodes(child)
	}
	return count
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		headerHeight := 3
		footerHeight := 4
		if m.showHelp {
			footerHeight = 15
		}
		verticalMarginHeight := headerHeight + footerHeight

		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
		m.help.Width = msg.Width

		m.updateViewport()
		return m, nil

	case tea.KeyMsg:
		if m.filterMode || m.jsonpathMode || m.searchMode || m.gotoMode {
			switch {
			case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
				if m.jsonpathMode {
					m.applyJSONPathFilter()
				} else if m.searchMode {
					m.performSearch()
				} else if m.gotoMode {
					m.gotoPath()
				}
				m.filterMode = false
				m.jsonpathMode = false
				m.searchMode = false
				m.gotoMode = false
				m.updateViewNodes()
				m.updateViewport()
				return m, nil
			case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
				m.filterMode = false
				m.jsonpathMode = false
				m.searchMode = false
				m.gotoMode = false
				m.filter = ""
				m.updateViewNodes()
				m.updateViewport()
				return m, nil
			case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
				if len(m.filter) > 0 {
					m.filter = m.filter[:len(m.filter)-1]
				}
				return m, nil
			default:
				m.filter += msg.String()
				return m, nil
			}
		}

		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			return m, nil
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.viewNodes)-1 {
				m.cursor++
				m.updateViewport()
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
				m.updateViewport()
			}
		case key.Matches(msg, m.keys.PageDown):
			m.cursor = min(len(m.viewNodes)-1, m.cursor+m.viewport.Height/2)
			m.updateViewport()
		case key.Matches(msg, m.keys.PageUp):
			m.cursor = max(0, m.cursor-m.viewport.Height/2)
			m.updateViewport()
		case key.Matches(msg, m.keys.Home):
			m.cursor = 0
			m.updateViewport()
		case key.Matches(msg, m.keys.End):
			m.cursor = len(m.viewNodes) - 1
			m.updateViewport()
		case key.Matches(msg, m.keys.Left):
			if m.cursor < len(m.viewNodes) {
				node := m.viewNodes[m.cursor]
				if node.Expanded && len(node.Children) > 0 {
					node.Expanded = false
					m.updateViewNodes()
					m.updateViewport()
				} else if node.Parent != nil {
					for i, n := range m.viewNodes {
						if n == node.Parent {
							m.cursor = i
							m.updateViewport()
							break
						}
					}
				}
			}
		case key.Matches(msg, m.keys.Right), key.Matches(msg, m.keys.Enter):
			if m.cursor < len(m.viewNodes) {
				node := m.viewNodes[m.cursor]
				if len(node.Children) > 0 {
					node.Expanded = !node.Expanded
					m.updateViewNodes()
					m.updateViewport()
				}
			}
		case key.Matches(msg, m.keys.ExpandAll):
			m.expandAll(m.root)
			m.updateViewNodes()
			m.updateViewport()
		case key.Matches(msg, m.keys.CollapseAll):
			m.collapseAll(m.root)
			m.updateViewNodes()
			m.updateViewport()
		case key.Matches(msg, m.keys.Reset):
			m.resetView()
		case key.Matches(msg, m.keys.Filter):
			m.filterMode = true
			m.filter = ""
		case key.Matches(msg, m.keys.JSONPath):
			m.jsonpathMode = true
			m.filter = ""
		case key.Matches(msg, m.keys.Search):
			m.searchMode = true
			m.filter = ""
		case key.Matches(msg, m.keys.Goto):
			m.gotoMode = true
			m.filter = ""
		case key.Matches(msg, m.keys.NextMatch):
			m.nextSearchMatch()
		case key.Matches(msg, m.keys.PrevMatch):
			m.prevSearchMatch()
		case key.Matches(msg, m.keys.Copy):
			if m.cursor < len(m.viewNodes) {
				node := m.viewNodes[m.cursor]
				var value string
				switch node.Type {
				case ObjectNode, ArrayNode:
					jsonBytes, _ := json.MarshalIndent(node.Value, "", "  ")
					value = string(jsonBytes)
				default:
					value = fmt.Sprintf("%v", node.Value)
				}
				clipboard.WriteAll(value)
			}
		case key.Matches(msg, m.keys.CopyPath):
			if m.cursor < len(m.viewNodes) {
				node := m.viewNodes[m.cursor]
				clipboard.WriteAll(node.Path)
			}
		case key.Matches(msg, m.keys.CopyKey):
			if m.cursor < len(m.viewNodes) {
				node := m.viewNodes[m.cursor]
				clipboard.WriteAll(node.Key)
			}
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	header := m.renderHeader()
	footer := m.renderFooter()
	
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
		footer,
	)
}

func (m Model) renderHeader() string {
	title := headerStyle.Render(fmt.Sprintf("JSON Viewer - %s", m.filename))
	
	stats := statusStyle.Render(fmt.Sprintf("Size: %.1fKB | Nodes: %d", 
		float64(m.fileSize)/1024, m.nodeCount))
	
	var filterInfo string
	if m.filterMode {
		filterInfo = filterStyle.Render(fmt.Sprintf("Filter: %s█", m.filter))
	} else if m.jsonpathMode {
		filterInfo = jsonpathStyle.Render(fmt.Sprintf("JSONPath: %s█", m.filter))
	} else if m.searchMode {
		filterInfo = searchStyle.Render(fmt.Sprintf("Search: %s█", m.filter))
	} else if m.gotoMode {
		filterInfo = gotoStyle.Render(fmt.Sprintf("Goto: %s█", m.filter))
	} else if m.filter != "" {
		if strings.HasPrefix(m.filter, "$") {
			filterInfo = jsonpathStyle.Render(fmt.Sprintf("Active JSONPath: %s", m.filter))
		} else {
			filterInfo = filterStyle.Render(fmt.Sprintf("Active Filter: %s", m.filter))
		}
	}

	breadcrumb := ""
	if m.cursor < len(m.viewNodes) && len(m.viewNodes) > 0 {
		node := m.viewNodes[m.cursor]
		breadcrumb = breadcrumbStyle.Render(fmt.Sprintf("Path: %s", node.Path))
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

func (m Model) renderFooter() string {
	if m.showHelp {
		return m.renderManualHelp()
	}

	if m.filterMode {
		return "\n" + filterStyle.Render("Press Enter to apply filter, Esc to cancel")
	} else if m.jsonpathMode {
		return "\n" + jsonpathStyle.Render("Press Enter to apply JSONPath, Esc to cancel")
	} else if m.searchMode {
		return "\n" + searchStyle.Render("Press Enter to search, Esc to cancel")
	} else if m.gotoMode {
		return "\n" + gotoStyle.Render("Press Enter to goto path, Esc to cancel")
	}

	helpText := statusStyle.Render("Press ? for help")
	return "\n" + helpText
}

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

func (m Model) renderNode(node *Node, isCursor bool) string {
	depth := m.getDepth(node)
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
		keyPart = keyStyle.Render(node.Key) + ": "
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
		valuePart = stringStyle.Render(fmt.Sprintf("\"%s\"", node.Value))
	case NumberNode:
		valuePart = numberStyle.Render(fmt.Sprintf("%v", node.Value))
	case BoolNode:
		valuePart = boolStyle.Render(fmt.Sprintf("%v", node.Value))
	case NullNode:
		valuePart = nullStyle.Render("null")
	}

	line := indent + icon + keyPart + valuePart

	if isCursor {
		line = cursorStyle.Render(line)
	}

	return line
}

func (m Model) getDepth(node *Node) int {
	depth := 0
	current := node.Parent
	for current != nil {
		depth++
		current = current.Parent
	}
	return depth
}

func (m *Model) updateViewNodes() {
	m.viewNodes = nil
	m.collectViewNodes(m.root)
	
	if m.filter != "" && !m.jsonpathMode {
		filtered := make([]*Node, 0)
		for _, node := range m.viewNodes {
			if strings.Contains(strings.ToLower(node.Path), strings.ToLower(m.filter)) ||
				strings.Contains(strings.ToLower(node.Key), strings.ToLower(m.filter)) {
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

func (m *Model) collectViewNodes(node *Node) {
	m.viewNodes = append(m.viewNodes, node)
	if node.Expanded {
		for _, child := range node.Children {
			m.collectViewNodes(child)
		}
	}
}

func (m *Model) applyJSONPathFilter() {
	if m.filter == "" {
		return
	}

	result, err := jsonpath.Get(m.filter, m.rawData)
	if err != nil {
		return
	}

	m.root = buildTree(result, "", "$")
	m.root.Expanded = true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *Model) expandAll(node *Node) {
	node.Expanded = true
	for _, child := range node.Children {
		m.expandAll(child)
	}
}

func (m *Model) collapseAll(node *Node) {
	node.Expanded = false
	for _, child := range node.Children {
		m.collapseAll(child)
	}
}

func (m *Model) resetView() {
	m.filter = ""
	m.filterMode = false
	m.jsonpathMode = false
	m.searchMode = false
	m.gotoMode = false
	m.searchMatches = nil
	m.searchIndex = 0
	m.cursor = 0
	
	m.root = buildTree(m.rawData, "", "$")
	m.root.Expanded = true
	m.updateViewNodes()
	m.updateViewport()
}

func (m *Model) performSearch() {
	if m.filter == "" {
		return
	}
	
	m.searchMatches = nil
	m.searchIndex = 0
	
	for _, node := range m.viewNodes {
		if strings.Contains(strings.ToLower(node.Key), strings.ToLower(m.filter)) ||
			strings.Contains(strings.ToLower(node.Path), strings.ToLower(m.filter)) ||
			(node.Value != nil && strings.Contains(strings.ToLower(fmt.Sprintf("%v", node.Value)), strings.ToLower(m.filter))) {
			m.searchMatches = append(m.searchMatches, node)
		}
	}
	
	if len(m.searchMatches) > 0 {
		m.jumpToSearchMatch(0)
	}
}

func (m *Model) nextSearchMatch() {
	if len(m.searchMatches) == 0 {
		return
	}
	
	m.searchIndex = (m.searchIndex + 1) % len(m.searchMatches)
	m.jumpToSearchMatch(m.searchIndex)
}

func (m *Model) prevSearchMatch() {
	if len(m.searchMatches) == 0 {
		return
	}
	
	m.searchIndex = (m.searchIndex - 1 + len(m.searchMatches)) % len(m.searchMatches)
	m.jumpToSearchMatch(m.searchIndex)
}

func (m *Model) jumpToSearchMatch(index int) {
	if index >= len(m.searchMatches) {
		return
	}
	
	targetNode := m.searchMatches[index]
	for i, node := range m.viewNodes {
		if node == targetNode {
			m.cursor = i
			m.updateViewport()
			break
		}
	}
}

func (m *Model) gotoPath() {
	if m.filter == "" {
		return
	}
	
	for i, node := range m.viewNodes {
		if strings.EqualFold(node.Path, m.filter) || strings.Contains(strings.ToLower(node.Path), strings.ToLower(m.filter)) {
			m.cursor = i
			m.updateViewport()
			break
		}
	}
}

func (m Model) renderManualHelp() string {
	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
	titleStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	
	help := strings.Builder{}
	help.WriteString("\n")
	help.WriteString(titleStyle.Render("━━━ JSON Viewer Help ━━━") + "\n\n")
	
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
	
	help.WriteString(helpStyle.Render("Clipboard:") + "\n")
	help.WriteString("  c, p, y                 Copy value/path/key\n\n")
	
	help.WriteString(helpStyle.Render("Utility:") + "\n")
	help.WriteString("  r/Ctrl+R                Reset view\n")
	help.WriteString("  ?, q/Esc                Help, Quit\n\n")
	
	help.WriteString(titleStyle.Render("Press ? to close help"))
	
	return help.String()
}