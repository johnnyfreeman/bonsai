package viewer

import "github.com/charmbracelet/lipgloss"

// DefaultTheme returns the default dark theme
func DefaultTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("86")),
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("208")),
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("196")),
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("240")),
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("171")).Bold(true),
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true),
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true),
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("220")).Foreground(lipgloss.Color("16")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("62")),
	}
}

// LightTheme returns a light theme
func LightTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("28")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("28")),
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("130")),
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("166")),
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("160")),
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("245")),
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("254")),
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("161")).Bold(true),
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("133")).Bold(true),
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Bold(true),
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("172")).Bold(true),
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("242")),
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("226")).Foreground(lipgloss.Color("16")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")),
	}
}

// MonochromeTheme returns a monochrome theme
func MonochromeTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("250")),
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("248")),
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("246")),
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("240")),
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Bold(true),
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("255")).Foreground(lipgloss.Color("16")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("240")),
	}
}