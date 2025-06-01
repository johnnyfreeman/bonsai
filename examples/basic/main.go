package main

import (
	"encoding/json"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnnyfreeman/bonsai/viewer"
)

func main() {
	// Sample JSON data
	jsonData := `{
		"name": "Example Application",
		"version": "1.0.0",
		"users": [
			{
				"id": 1,
				"name": "Alice",
				"active": true,
				"profile": {
					"email": "alice@example.com",
					"age": 30
				}
			},
			{
				"id": 2,
				"name": "Bob",
				"active": false,
				"profile": {
					"email": "bob@example.com",
					"age": 25
				}
			}
		],
		"config": {
			"debug": true,
			"max_connections": 100,
			"timeout": null
		}
	}`

	// Parse JSON
	var data interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal(err)
	}

	// Create viewer with default configuration
	model := viewer.New(data)

	// Run the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}