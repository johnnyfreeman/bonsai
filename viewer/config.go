package viewer

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Theme:             DefaultTheme(),
		ShowHelp:          true,
		ShowLineNumbers:   false,
		ShowBorders:       true,
		InitiallyExpanded: true,
		EnableMouse:       false,
		EnableClipboard:   true,
		Width:             0, // 0 means auto-size
		Height:            0, // 0 means auto-size
		
		// Default no-op callbacks
		OnSelect:   func(*Node) {},
		OnExpand:   func(*Node) {},
		OnCollapse: func(*Node) {},
		OnCopy:     func(string) {},
		OnFilter:   func(string) {},
		OnError:    func(error) {},
	}
}

// WithTheme sets the theme
func (c Config) WithTheme(theme Theme) Config {
	c.Theme = theme
	return c
}

// WithSize sets the size (for embedded mode)
func (c Config) WithSize(width, height int) Config {
	c.Width = width
	c.Height = height
	return c
}

// WithCallbacks sets the event callbacks
func (c Config) WithCallbacks(onSelect, onExpand, onCollapse func(*Node)) Config {
	if onSelect != nil {
		c.OnSelect = onSelect
	}
	if onExpand != nil {
		c.OnExpand = onExpand
	}
	if onCollapse != nil {
		c.OnCollapse = onCollapse
	}
	return c
}

// WithClipboard sets clipboard callback
func (c Config) WithClipboard(onCopy func(string)) Config {
	if onCopy != nil {
		c.OnCopy = onCopy
	}
	return c
}

// WithFilter sets filter callback
func (c Config) WithFilter(onFilter func(string)) Config {
	if onFilter != nil {
		c.OnFilter = onFilter
	}
	return c
}

// WithError sets error callback
func (c Config) WithError(onError func(error)) Config {
	if onError != nil {
		c.OnError = onError
	}
	return c
}

// Embedded returns a config optimized for embedding
func (c Config) Embedded() Config {
	c.ShowHelp = false
	c.ShowBorders = false
	return c
}

// ReadOnly returns a config for read-only viewing
func (c Config) ReadOnly() Config {
	c.EnableClipboard = false
	return c
}