# Bonsai

A powerful, terminal-based JSON viewer library and CLI tool with vim-like navigation, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **🚀 Library + CLI**: Use as a standalone tool or embed in your Bubble Tea applications
- **⌨️ Vim-like Navigation**: hjkl keys, page up/down, home/end
- **🔍 Powerful Filtering**: Text filters, JSONPath queries, content search
- **🎨 Customizable Themes**: Built-in themes (dark, light, monochrome) and full customization
- **📋 Clipboard Support**: Copy values, paths, or keys
- **🔧 Event System**: Callbacks for selections, expansions, and user actions
- **📱 Embedded Mode**: Perfect for integrating into larger applications
- **⚡ High Performance**: Efficiently handles large JSON files

## Installation

### As a CLI tool:
```bash
go install github.com/johnnyfreeman/bonsai/cmd@latest
```

### As a library:
```bash
go get github.com/johnnyfreeman/bonsai/viewer
```

## CLI Usage

```bash
bonsai example.json
```

### Controls

#### Navigation
- `hjkl` or arrow keys: Navigate through the JSON tree
- `PgUp`/`Ctrl+U`: Page up (half screen)
- `PgDn`/`Ctrl+D`: Page down (half screen)
- `Home`/`g`: Go to top
- `End`/`G`: Go to bottom

#### Tree Operations
- `Enter`/`Space`/`l`: Expand/collapse node
- `h`: Collapse current node or move to parent
- `E`: Expand all nodes
- `C`: Collapse all nodes

#### Filtering & Search
- `/`: Enter text filter mode
- `$`: Enter JSONPath filter mode
- `s`/`Ctrl+F`: Search in content
- `n`: Next search match
- `N`: Previous search match
- `:`/`Ctrl+G`: Goto path

#### Clipboard Operations
- `c`: Copy current value
- `p`: Copy current path
- `y`: Copy current key

#### Utility
- `r`/`Ctrl+R`: Reset view (clear filters)
- `?`: Toggle help
- `q`/`Esc`/`Ctrl+C`: Quit

## Library Usage

### Basic Usage

```go
package main

import (
    "encoding/json"
    "log"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/johnnyfreeman/bonsai/viewer"
)

func main() {
    // Parse your JSON data
    var data interface{}
    json.Unmarshal(jsonBytes, &data)
    
    // Create viewer with default config
    model := viewer.New(data)
    
    // Run as a Bubble Tea program
    tea.NewProgram(model, tea.WithAltScreen()).Run()
}
```

### Custom Configuration

```go
config := viewer.DefaultConfig().
    WithTheme(viewer.LightTheme()).
    WithCallbacks(
        func(node *viewer.Node) { /* on select */ },
        func(node *viewer.Node) { /* on expand */ },
        func(node *viewer.Node) { /* on collapse */ },
    ).
    WithClipboard(func(content string) {
        // Handle clipboard operations
    }).
    WithError(func(err error) {
        // Handle errors
    })

model := viewer.New(data, config)
```

### Embedded Mode

Perfect for integrating into larger applications:

```go
config := viewer.DefaultConfig().
    Embedded().                    // Remove borders and help
    WithSize(60, 20).             // Fixed size
    ReadOnly()                    // Disable clipboard

// Embed in your larger Bubble Tea model
type MyApp struct {
    jsonViewer viewer.Model
    // ... other components
}
```

### Available Themes

```go
// Built-in themes
viewer.DefaultTheme()         // Dark theme
viewer.LightTheme()          // Light theme  
viewer.MonochromeTheme()     // Black and white

// Popular community themes
viewer.TokyoNightTheme()     // Popular VS Code theme
viewer.CatppuccinMochaTheme() // Catppuccin dark variant
viewer.CatppuccinLatteTheme() // Catppuccin light variant
viewer.DraculaTheme()        // Classic dark theme
viewer.NordTheme()           // Arctic/minimal theme
viewer.GruvboxTheme()        // Retro warm theme

// Or create your own
customTheme := viewer.Theme{
    Header: lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B6B")),
    Key:    lipgloss.NewStyle().Foreground(lipgloss.Color("#4ECDC4")),
    // ... customize all elements
}
```

### Event Callbacks

```go
config := viewer.DefaultConfig().
    WithCallbacks(
        // OnSelect - fired when cursor moves
        func(node *viewer.Node) {
            fmt.Printf("Selected: %s\n", node.Path)
        },
        // OnExpand - fired when node is expanded
        func(node *viewer.Node) {
            fmt.Printf("Expanded: %s\n", node.Path)  
        },
        // OnCollapse - fired when node is collapsed
        func(node *viewer.Node) {
            fmt.Printf("Collapsed: %s\n", node.Path)
        },
    ).
    WithFilter(func(filter string) {
        fmt.Printf("Filter applied: %s\n", filter)
    }).
    WithClipboard(func(content string) {
        fmt.Printf("Copied: %s\n", content)
    })
```

## API Reference

### Core Types

```go
// Main model
type Model struct { ... }

// Node represents a JSON node
type Node struct {
    Key      string
    Value    interface{}
    Type     NodeType
    Children []*Node
    Parent   *Node
    Expanded bool
    Path     string
}

// Configuration
type Config struct {
    Theme           Theme
    ShowHelp        bool
    ShowBorders     bool
    EnableClipboard bool
    // ... callbacks and options
}
```

### Key Methods

```go
// Creation
viewer.New(data interface{}, config ...Config) Model
viewer.NewFromJSON([]byte, config ...Config) (Model, error)
viewer.NewFromReader(io.Reader, config ...Config) (Model, error)

// Querying
model.GetCurrentNode() *Node
model.GetFilteredData() interface{}
model.IsFiltered() bool
model.GetSearchMatches() []*Node

// Configuration
config.WithTheme(Theme) Config
config.WithSize(width, height int) Config
config.Embedded() Config
config.ReadOnly() Config
```

## Examples

Check out the [examples/](examples/) directory for complete examples:

- **[basic/](examples/basic/)**: Simple standalone viewer
- **[embedded/](examples/embedded/)**: Embedding in a larger application  
- **[custom-theme/](examples/custom-theme/)**: Custom styling and callbacks

## JSONPath Examples

Press `$` to enter JSONPath mode with smart path suggestions:

**Smart Path Suggestions:** When you press `$`, the viewer intelligently suggests a JSONPath based on your current cursor position. For example, if you're on `$.users[0].name`, it will start with `$.users[*].name` to show all user names.

- `$.users[*].name` - Get all user names
- `$.config.database` - Get database config
- `$.items[?(@.active == true)]` - Get active items
- `$..email` - Get all email fields recursively

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.