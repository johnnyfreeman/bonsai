package main

import (
	"encoding/json"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/johnnyfreeman/bonsai/viewer"
)

// App demonstrates embedding the JSON viewer in a larger application
type App struct {
	jsonViewer viewer.Model
	sidebar    string
	logs       []string
}

func (a App) Init() tea.Cmd {
	return a.jsonViewer.Init()
}

func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return a, tea.Quit
		}
	}

	// Update the embedded JSON viewer
	a.jsonViewer, cmd = a.jsonViewer.Update(msg)

	return a, cmd
}

func (a App) View() string {
	// Create a 3-column layout
	sidebarStyle := lipgloss.NewStyle().
		Width(20).
		Height(20).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1)

	viewerStyle := lipgloss.NewStyle().
		Width(60).
		Height(20).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86"))

	logsStyle := lipgloss.NewStyle().
		Width(30).
		Height(20).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("214")).
		Padding(1)

	sidebar := sidebarStyle.Render(a.sidebar)
	jsonView := viewerStyle.Render(a.jsonViewer.View())
	logs := logsStyle.Render(strings.Join(a.logs, "\n"))

	return lipgloss.JoinHorizontal(lipgloss.Top, sidebar, jsonView, logs)
}

func main() {
	// Sample data
	jsonData := `{
		"application": {
			"name": "Embedded Demo",
			"components": ["viewer", "sidebar", "logs"],
			"settings": {
				"theme": "dark",
				"layout": "horizontal"
			}
		}
	}`

	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal(err)
	}

	// Create embedded viewer configuration
	config := viewer.DefaultConfig().
		Embedded().
		WithSize(58, 18). // Fit within the styled border
		WithCallbacks(
			func(node *viewer.Node) {
				// Handle selections - could update sidebar or logs
			},
			nil, nil,
		)

	// Create the embedded viewer
	jsonViewer := viewer.New(data, config)

	// Create the main app
	app := App{
		jsonViewer: jsonViewer,
		sidebar: `Navigation:
• File Browser
• Settings
• Tools

Status:
• JSON loaded
• 7 nodes
• Ready`,
		logs: []string{
			"[INFO] App started",
			"[INFO] JSON loaded",
			"[INFO] Viewer ready",
			"[DEBUG] Embedded mode",
		},
	}

	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}