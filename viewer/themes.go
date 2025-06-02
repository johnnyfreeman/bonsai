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

// TokyoNightTheme returns the popular Tokyo Night theme
func TokyoNightTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true), // Blue
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")),           // Comment
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("#7dcfff")),           // Cyan
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("#9ece6a")),           // Green
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9e64")),           // Orange
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("#f7768e")),           // Red
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")),           // Comment
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("#283457")),           // Selection
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("#bb9af7")).Bold(true), // Purple
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true), // Blue
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("#9ece6a")).Bold(true), // Green
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("#e0af68")).Bold(true), // Yellow
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("#9aa5ce")),           // Fg dark
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("#e0af68")).Foreground(lipgloss.Color("#1a1b26")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#414868")),
	}
}

// CatppuccinMochaTheme returns the Catppuccin Mocha (dark) theme
func CatppuccinMochaTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa")).Bold(true), // Blue
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")),           // Overlay1
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("#94e2d5")),           // Teal
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")),           // Green
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("#fab387")),           // Peach
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("#f38ba8")),           // Pink
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")),           // Overlay1
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("#313244")),           // Surface0
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7")).Bold(true), // Mauve
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("#89b4fa")).Bold(true), // Blue
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("#a6e3a1")).Bold(true), // Green
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("#f9e2af")).Bold(true), // Yellow
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("#bac2de")),           // Subtext1
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("#f9e2af")).Foreground(lipgloss.Color("#1e1e2e")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#45475a")),
	}
}

// CatppuccinLatteTheme returns the Catppuccin Latte (light) theme
func CatppuccinLatteTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("#1e66f5")).Bold(true), // Blue
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#8c8fa1")),           // Overlay1
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("#179299")),           // Teal
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("#40a02b")),           // Green
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("#fe640b")),           // Peach
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ea76cb")),           // Pink
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("#8c8fa1")),           // Overlay1
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("#e6e9ef")),           // Surface0
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("#8839ef")).Bold(true), // Mauve
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("#1e66f5")).Bold(true), // Blue
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("#40a02b")).Bold(true), // Green
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("#df8e1d")).Bold(true), // Yellow
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("#6c6f85")),           // Subtext1
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("#df8e1d")).Foreground(lipgloss.Color("#eff1f5")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#9ca0b0")),
	}
}

// DraculaTheme returns the classic Dracula theme
func DraculaTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Bold(true), // Purple
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4")),           // Comment
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")),           // Cyan
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")),           // Green
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb86c")),           // Orange
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ff79c6")),           // Pink
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4")),           // Comment
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("#44475a")),           // Selection
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Bold(true), // Purple
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("#8be9fd")).Bold(true), // Cyan
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Bold(true), // Green
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c")).Bold(true), // Yellow
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("#f8f8f2")),           // Foreground
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("#f1fa8c")).Foreground(lipgloss.Color("#282a36")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#6272a4")),
	}
}

// NordTheme returns the Nord arctic theme
func NordTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("#81a1c1")).Bold(true), // Nord9
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#4c566a")),           // Nord3
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("#88c0d0")),           // Nord8
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("#a3be8c")),           // Nord14
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("#d08770")),           // Nord12
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("#bf616a")),           // Nord11
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("#4c566a")),           // Nord3
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("#3b4252")),           // Nord1
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("#b48ead")).Bold(true), // Nord15
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("#5e81ac")).Bold(true), // Nord10
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("#a3be8c")).Bold(true), // Nord14
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("#ebcb8b")).Bold(true), // Nord13
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("#d8dee9")),           // Nord4
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("#ebcb8b")).Foreground(lipgloss.Color("#2e3440")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#434c5e")),
	}
}

// GruvboxTheme returns the Gruvbox retro theme
func GruvboxTheme() Theme {
	return Theme{
		Header:     lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598")).Bold(true), // Bright blue
		Status:     lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")),           // Gray
		Key:        lipgloss.NewStyle().Foreground(lipgloss.Color("#8ec07c")),           // Bright aqua
		String:     lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")),           // Bright green
		Number:     lipgloss.NewStyle().Foreground(lipgloss.Color("#fe8019")),           // Bright orange
		Bool:       lipgloss.NewStyle().Foreground(lipgloss.Color("#fb4934")),           // Bright red
		Null:       lipgloss.NewStyle().Foreground(lipgloss.Color("#928374")),           // Gray
		Cursor:     lipgloss.NewStyle().Background(lipgloss.Color("#3c3836")),           // Dark1
		Filter:     lipgloss.NewStyle().Foreground(lipgloss.Color("#d3869b")).Bold(true), // Bright purple
		JSONPath:   lipgloss.NewStyle().Foreground(lipgloss.Color("#83a598")).Bold(true), // Bright blue
		Search:     lipgloss.NewStyle().Foreground(lipgloss.Color("#b8bb26")).Bold(true), // Bright green
		Goto:       lipgloss.NewStyle().Foreground(lipgloss.Color("#fabd2f")).Bold(true), // Bright yellow
		Breadcrumb: lipgloss.NewStyle().Foreground(lipgloss.Color("#a89984")),           // Light4
		Match:      lipgloss.NewStyle().Background(lipgloss.Color("#fabd2f")).Foreground(lipgloss.Color("#1d2021")),
		Border:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#504945")),
	}
}