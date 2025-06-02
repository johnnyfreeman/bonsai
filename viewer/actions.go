package viewer

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

// Navigation actions
func (m *Model) moveDown() {
	if m.cursor < len(m.viewNodes)-1 {
		m.cursor++
		m.updateViewport()
		if node := m.GetCurrentNode(); node != nil {
			m.config.OnSelect(node)
		}
	}
}

func (m *Model) moveUp() {
	if m.cursor > 0 {
		m.cursor--
		m.updateViewport()
		if node := m.GetCurrentNode(); node != nil {
			m.config.OnSelect(node)
		}
	}
}

func (m *Model) pageDown() {
	m.cursor = min(len(m.viewNodes)-1, m.cursor+m.viewport.Height/2)
	m.updateViewport()
	if node := m.GetCurrentNode(); node != nil {
		m.config.OnSelect(node)
	}
}

func (m *Model) pageUp() {
	m.cursor = max(0, m.cursor-m.viewport.Height/2)
	m.updateViewport()
	if node := m.GetCurrentNode(); node != nil {
		m.config.OnSelect(node)
	}
}

func (m *Model) goHome() {
	m.cursor = 0
	m.updateViewport()
	if node := m.GetCurrentNode(); node != nil {
		m.config.OnSelect(node)
	}
}

func (m *Model) goEnd() {
	m.cursor = len(m.viewNodes) - 1
	m.updateViewport()
	if node := m.GetCurrentNode(); node != nil {
		m.config.OnSelect(node)
	}
}

func (m *Model) moveLeft() {
	if m.cursor < len(m.viewNodes) {
		node := m.viewNodes[m.cursor]
		if node.Expanded && len(node.Children) > 0 {
			node.Expanded = false
			m.config.OnCollapse(node)
			m.updateViewNodes()
			m.updateViewport()
		} else if node.Parent != nil {
			for i, n := range m.viewNodes {
				if n == node.Parent {
					m.cursor = i
					m.updateViewport()
					m.config.OnSelect(n)
					break
				}
			}
		}
	}
}

func (m *Model) toggleExpansion() {
	if m.cursor < len(m.viewNodes) {
		node := m.viewNodes[m.cursor]
		if len(node.Children) > 0 {
			node.Expanded = !node.Expanded
			if node.Expanded {
				m.config.OnExpand(node)
			} else {
				m.config.OnCollapse(node)
			}
			m.updateViewNodes()
			m.updateViewport()
		}
	}
}

// Tree operations
func (m *Model) expandAll() {
	m.root.ExpandAll()
	m.updateViewNodes()
	m.updateViewport()
}

func (m *Model) collapseAll() {
	m.root.CollapseAll()
	m.updateViewNodes()
	m.updateViewport()
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

	m.root = BuildTree(m.rawData, "", "$")
	if m.config.InitiallyExpanded {
		m.root.Expanded = true
	}
	m.updateViewNodes()
	m.updateViewport()
}

// Input mode actions
func (m *Model) enterFilterMode() {
	m.filterMode = true
	m.filter = ""
}

func (m *Model) enterJSONPathMode() {
	m.jsonpathMode = true
	m.filter = "$"
	// Apply the initial filter to show the live filtering immediately
	m.applyLiveJSONPathFilter()
}

func (m *Model) enterSearchMode() {
	m.searchMode = true
	m.filter = ""
}

func (m *Model) enterGotoMode() {
	m.gotoMode = true
	m.filter = ""
}

// Input handling
func (m Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
		m.applyInput()
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("esc"))):
		m.cancelInput()
		return m, nil
	case key.Matches(msg, key.NewBinding(key.WithKeys("backspace"))):
		if len(m.filter) > 0 {
			m.filter = m.filter[:len(m.filter)-1]
			// Apply live filtering for text filter mode
			if m.filterMode {
				m.applyLiveFilter()
			} else if m.jsonpathMode {
				m.applyLiveJSONPathFilter()
			}
		}
		return m, nil
	default:
		m.filter += msg.String()
		// Apply live filtering for text filter mode
		if m.filterMode {
			m.applyLiveFilter()
		} else if m.jsonpathMode {
			m.applyLiveJSONPathFilter()
		}
		return m, nil
	}
}

func (m *Model) applyInput() {
	wasJSONPathMode := m.jsonpathMode
	
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
	
	if m.filter != "" {
		m.config.OnFilter(m.filter)
	}
	
	// Only update view nodes if it wasn't JSONPath mode (JSONPath already updates the tree structure)
	if !wasJSONPathMode {
		m.updateViewNodes()
		m.updateViewport()
	}
}

func (m *Model) cancelInput() {
	m.filterMode = false
	m.jsonpathMode = false
	m.searchMode = false
	m.gotoMode = false
	m.filter = ""
	m.updateViewNodes()
	m.updateViewport()
}

// applyLiveFilter applies text filtering as the user types
func (m *Model) applyLiveFilter() {
	// Save current cursor position to try to preserve it
	var currentNode *Node
	if m.cursor >= 0 && m.cursor < len(m.viewNodes) {
		currentNode = m.viewNodes[m.cursor]
	}
	
	// Update view nodes with current filter
	m.updateViewNodes()
	
	// Try to preserve cursor position by finding the same node
	if currentNode != nil {
		for i, node := range m.viewNodes {
			if node == currentNode {
				m.cursor = i
				break
			}
		}
		// If node not found, keep cursor at 0 or adjust to valid range
		if m.cursor >= len(m.viewNodes) {
			m.cursor = max(0, len(m.viewNodes)-1)
		}
	}
	
	m.updateViewport()
}

// applyLiveJSONPathFilter applies JSONPath filtering as the user types
func (m *Model) applyLiveJSONPathFilter() {
	// Don't apply empty filters
	if m.filter == "" {
		// Reset to original data when filter is empty
		m.root = BuildTree(m.rawData, "", "$")
		if m.config.InitiallyExpanded {
			m.root.Expanded = true
		}
		m.cursor = 0
		m.updateViewNodes()
		m.updateViewport()
		return
	}

	// Try to apply JSONPath filter, but don't show errors during live typing
	// Only apply if it's a potentially valid JSONPath (starts with $ or has some basic structure)
	if strings.HasPrefix(m.filter, "$") || strings.Contains(m.filter, ".") {
		result, err := jsonpath.Get(m.filter, m.rawData)
		if err == nil {
			m.root = BuildTree(result, "", "$")
			m.root.Expanded = true
			
			// Reset cursor to top since structure changed significantly
			m.cursor = 0
			
			m.updateViewNodes()
			m.updateViewport()
		}
		// Silently ignore errors during live typing - they'll be shown on Enter
	}
}

// Filtering and search
func (m *Model) applyJSONPathFilter() {
	if m.filter == "" {
		return
	}

	result, err := jsonpath.Get(m.filter, m.rawData)
	if err != nil {
		m.config.OnError(err)
		return
	}

	m.root = BuildTree(result, "", "$")
	m.root.Expanded = true
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
			m.config.OnSelect(node)
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
			m.config.OnSelect(node)
			break
		}
	}
}

// Clipboard operations
func (m *Model) copyValue() {
	if !m.config.EnableClipboard {
		return
	}

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
		m.config.OnCopy(value)
	}
}

func (m *Model) copyPath() {
	if !m.config.EnableClipboard {
		return
	}

	if m.cursor < len(m.viewNodes) {
		node := m.viewNodes[m.cursor]
		clipboard.WriteAll(node.Path)
		m.config.OnCopy(node.Path)
	}
}

func (m *Model) copyKey() {
	if !m.config.EnableClipboard {
		return
	}

	if m.cursor < len(m.viewNodes) {
		node := m.viewNodes[m.cursor]
		clipboard.WriteAll(node.Key)
		m.config.OnCopy(node.Key)
	}
}

// Utility functions
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