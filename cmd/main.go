package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/johnnyfreeman/bonsai/viewer"
)

func main() {
	var data []byte
	var filename string
	var fileSize int64
	var err error

	// Check if we have stdin data or a file argument
	stat, _ := os.Stdin.Stat()
	hasStdin := (stat.Mode() & os.ModeCharDevice) == 0

	if hasStdin {
		// Read from stdin
		data, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Error reading from stdin: %v", err)
		}
		filename = "<stdin>"
		fileSize = int64(len(data))
	} else {
		// Check for file argument
		if len(os.Args) < 2 {
			fmt.Println("Usage: bonsai <file.json>")
			fmt.Println("   or: cat file.json | bonsai")
			fmt.Println("   or: curl -s api.example.com/data.json | bonsai")
			fmt.Println("\nBonsai - A terminal-based JSON viewer with vim-like navigation.")
			fmt.Println("\nFeatures:")
			fmt.Println("  • hjkl navigation")
			fmt.Println("  • Expand/collapse nodes")
			fmt.Println("  • Text and JSONPath filtering")
			fmt.Println("  • Search functionality")
			fmt.Println("  • Copy to clipboard")
			fmt.Println("  • Multiple themes")
			fmt.Println("\nPress ? for help when running")
			os.Exit(1)
		}

		filename = os.Args[1]

		// Read the file
		file, err := os.Open(filename)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer file.Close()

		// Get file info for display
		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatalf("Error getting file info: %v", err)
		}
		fileSize = fileInfo.Size()

		// Read file content
		data, err = io.ReadAll(file)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
	}

	// Create the viewer model with CLI-specific configuration
	config := viewer.DefaultConfig().
		WithTheme(viewer.TokyoNightTheme())

	// Set up error handling
	config.OnError = func(err error) {
		log.Printf("Bonsai Error: %v", err)
	}

	// Create the model
	model, err := viewer.NewFromJSON(data, config)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Add file information
	if filename == "<stdin>" {
		model = model.WithFilename(filename, fileSize)
	} else {
		model = model.WithFilename(filepath.Base(filename), fileSize)
	}

	// Run the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}
}
