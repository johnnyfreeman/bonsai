package main

import (
	"encoding/json"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johnnyfreeman/bonsai/viewer"
)

func main() {
	// Sample JSON data
	jsonData := `{
		"theme": "custom",
		"colors": {
			"primary": "#FF6B6B",
			"secondary": "#4ECDC4",
			"accent": "#45B7D1"
		},
		"data": [
			{"name": "Red", "hex": "#FF0000", "rgb": [255, 0, 0]},
			{"name": "Green", "hex": "#00FF00", "rgb": [0, 255, 0]},
			{"name": "Blue", "hex": "#0000FF", "rgb": [0, 0, 255]}
		],
		"settings": {
			"animated": true,
			"duration": 300,
			"easing": "ease-in-out"
		}
	}`

	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal(err)
	}

	// Create a custom theme with vibrant colors
	customTheme := viewer.Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")),
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4")).Bold(true),
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FFE66D")),
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF8E53")),
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("#95E1D3")),
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")),
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("#45B7D1")).Foreground(lipgloss.Color("#000000")),
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")).Bold(true),
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("#C44569")).Bold(true),
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("#26de81")).Bold(true),
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("#fd9644")).Bold(true),
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("#a55eea")),
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("#FFE66D")).Foreground(lipgloss.Color("#000000")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#4ECDC4")),
	}

	// Create configuration with custom theme and callbacks
	config := viewer.DefaultConfig().
		WithTheme(customTheme).
		WithCallbacks(
			func(node *viewer.Node) {
				// Log selections (in a real app, you might update other UI elements)
				log.Printf("Selected: %s = %v", node.Path, node.Value)
			},
			func(node *viewer.Node) {
				log.Printf("Expanded: %s", node.Path)
			},
			func(node *viewer.Node) {
				log.Printf("Collapsed: %s", node.Path)
			},
		).
		WithClipboard(func(content string) {
			log.Printf("Copied to clipboard: %s", content)
		}).
		WithFilter(func(filter string) {
			log.Printf("Filter applied: %s", filter)
		}).
		WithError(func(err error) {
			log.Printf("Error: %v", err)
		})

	// Create the viewer
	model := viewer.New(data, config)

	// Run the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}