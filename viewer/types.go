package viewer

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// NodeType represents the type of a JSON node
type NodeType int

const (
	ObjectNode NodeType = iota
	ArrayNode
	StringNode
	NumberNode
	BoolNode
	NullNode
)

// Node represents a node in the JSON tree
type Node struct {
	Key      string
	Value    interface{}
	Type     NodeType
	Children []*Node
	Parent   *Node
	Expanded bool
	Path     string
}

// Config holds configuration options for the JSON viewer
type Config struct {
	// Appearance
	Theme          Theme
	ShowHelp       bool
	ShowLineNumbers bool
	ShowBorders    bool
	
	// Behavior
	InitiallyExpanded bool
	EnableMouse       bool
	EnableClipboard   bool
	
	// Callbacks
	OnSelect     func(*Node)
	OnExpand     func(*Node)
	OnCollapse   func(*Node)
	OnCopy       func(string)
	OnFilter     func(string)
	OnError      func(error)
	
	// Size constraints (for embedded mode)
	Width  int
	Height int
}

// Theme defines the color scheme and styling
type Theme struct {
	Header      lipgloss.Style
	Status      lipgloss.Style
	Key         lipgloss.Style
	String      lipgloss.Style
	Number      lipgloss.Style
	Bool        lipgloss.Style
	Null        lipgloss.Style
	Cursor      lipgloss.Style
	Filter      lipgloss.Style
	JSONPath    lipgloss.Style
	Search      lipgloss.Style
	Goto        lipgloss.Style
	Breadcrumb  lipgloss.Style
	Match       lipgloss.Style
	Border      lipgloss.Style
}

// KeyMap defines the key bindings for the viewer
type KeyMap struct {
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

// Model is the Bubble Tea model for the JSON viewer
type Model struct {
	// Core data
	root          *Node
	rawData       interface{}
	config        Config
	
	// UI state
	cursor        int
	viewNodes     []*Node
	viewport      viewport.Model
	help          help.Model
	progress      progress.Model
	keys          KeyMap
	
	// File info
	filename      string
	fileSize      int64
	nodeCount     int
	
	// Mode state
	filter        string
	filterMode    bool
	jsonpathMode  bool
	searchMode    bool
	gotoMode      bool
	showHelp      bool
	
	// Search state
	searchMatches []*Node
	searchIndex   int
	
	// Position restoration
	savedNodePath string
	
	// Embedded mode
	embedded      bool
	width         int
	height        int
}