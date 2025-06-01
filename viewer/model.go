package viewer

import (
	"encoding/json"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// New creates a new JSON viewer model
func New(data interface{}, config ...Config) Model {
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	root := BuildTree(data, "", "$")
	if cfg.InitiallyExpanded {
		root.Expanded = true
	}

	vp := viewport.New(80, 20)
	if cfg.ShowBorders {
		vp.Style = cfg.Theme.Border.PaddingRight(2)
	}

	helpModel := help.New()
	helpModel.Width = 80

	m := Model{
		root:      root,
		rawData:   data,
		config:    cfg,
		cursor:    0,
		viewport:  vp,
		help:      helpModel,
		progress:  progress.New(progress.WithDefaultGradient()),
		keys:      DefaultKeyMap(),
		nodeCount: CountNodes(root),
		embedded:  cfg.Width > 0 || cfg.Height > 0,
		width:     cfg.Width,
		height:    cfg.Height,
	}

	m.updateViewNodes()
	m.updateViewport()

	return m
}

// NewFromJSON creates a new JSON viewer from JSON data
func NewFromJSON(jsonData []byte, config ...Config) (Model, error) {
	var data interface{}
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return Model{}, err
	}
	return New(data, config...), nil
}

// NewFromReader creates a new JSON viewer from an io.Reader
func NewFromReader(reader io.Reader, config ...Config) (Model, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return Model{}, err
	}
	return NewFromJSON(data, config...)
}

// WithFilename sets the filename for display purposes
func (m Model) WithFilename(filename string, size int64) Model {
	m.filename = filename
	m.fileSize = size
	return m
}

// GetCurrentNode returns the currently selected node
func (m Model) GetCurrentNode() *Node {
	if m.cursor >= 0 && m.cursor < len(m.viewNodes) {
		return m.viewNodes[m.cursor]
	}
	return nil
}

// GetFilteredData returns the current filtered data
func (m Model) GetFilteredData() interface{} {
	if m.filter == "" {
		return m.rawData
	}
	// Return the current root's value which represents filtered data
	return m.root.Value
}

// IsFiltered returns true if any filter is currently active
func (m Model) IsFiltered() bool {
	return m.filter != ""
}

// GetSearchMatches returns the current search matches
func (m Model) GetSearchMatches() []*Node {
	return m.searchMatches
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.embedded {
			headerHeight := 3
			footerHeight := 4
			if m.showHelp {
				footerHeight = 15
			}
			verticalMarginHeight := headerHeight + footerHeight

			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMarginHeight
			m.help.Width = msg.Width
		} else {
			// Embedded mode uses fixed size or adapts to available space
			if m.width > 0 {
				m.viewport.Width = m.width
			} else {
				m.viewport.Width = msg.Width
			}
			if m.height > 0 {
				m.viewport.Height = m.height
			} else {
				m.viewport.Height = msg.Height
			}
		}

		m.updateViewport()
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}

	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// handleKeyPress handles key press events
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle input modes first
	if m.filterMode || m.jsonpathMode || m.searchMode || m.gotoMode {
		return m.handleInputMode(msg)
	}

	// Handle regular key presses
	switch {
	case key.Matches(msg, m.keys.Quit):
		return m, tea.Quit
	case key.Matches(msg, m.keys.Help):
		if m.config.ShowHelp {
			m.showHelp = !m.showHelp
		}
		return m, nil
	case key.Matches(msg, m.keys.Down):
		m.moveDown()
	case key.Matches(msg, m.keys.Up):
		m.moveUp()
	case key.Matches(msg, m.keys.PageDown):
		m.pageDown()
	case key.Matches(msg, m.keys.PageUp):
		m.pageUp()
	case key.Matches(msg, m.keys.Home):
		m.goHome()
	case key.Matches(msg, m.keys.End):
		m.goEnd()
	case key.Matches(msg, m.keys.Left):
		m.moveLeft()
	case key.Matches(msg, m.keys.Right), key.Matches(msg, m.keys.Enter):
		m.toggleExpansion()
	case key.Matches(msg, m.keys.ExpandAll):
		m.expandAll()
	case key.Matches(msg, m.keys.CollapseAll):
		m.collapseAll()
	case key.Matches(msg, m.keys.Reset):
		m.resetView()
	case key.Matches(msg, m.keys.Filter):
		m.enterFilterMode()
	case key.Matches(msg, m.keys.JSONPath):
		m.enterJSONPathMode()
	case key.Matches(msg, m.keys.Search):
		m.enterSearchMode()
	case key.Matches(msg, m.keys.Goto):
		m.enterGotoMode()
	case key.Matches(msg, m.keys.NextMatch):
		m.nextSearchMatch()
	case key.Matches(msg, m.keys.PrevMatch):
		m.prevSearchMatch()
	case key.Matches(msg, m.keys.Copy):
		m.copyValue()
	case key.Matches(msg, m.keys.CopyPath):
		m.copyPath()
	case key.Matches(msg, m.keys.CopyKey):
		m.copyKey()
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if m.embedded {
		return m.viewport.View()
	}

	header := m.renderHeader()
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		m.viewport.View(),
		footer,
	)
}