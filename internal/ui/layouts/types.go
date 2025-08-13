package layouts

import (
	"github.com/charmbracelet/lipgloss"
)

// Config holds the layout manager configuration
type Config struct {
	// Theme configuration
	Theme ThemeConfig

	// Layout behavior
	EnableModal      bool
	EnableStatusBar  bool
	EnableHelpBar    bool
	ResponsiveLayout bool

	// Spacing and margins
	ContentMargin  int
	ContentPadding int
}

// DefaultConfig returns a default layout configuration
func DefaultConfig() Config {
	return Config{
		Theme:            DefaultTheme(),
		EnableModal:      true,
		EnableStatusBar:  true,
		EnableHelpBar:    true,
		ResponsiveLayout: true,
		ContentMargin:    0,
		ContentPadding:   0,
	}
}

// LayoutOptions contains options for creating layouts
type LayoutOptions struct {
	// Content styling
	ContentStyle *lipgloss.Style

	// Context for components
	HelpContext   HelpContext
	StatusContext StatusContext

	// Modal overlay
	Modal *ModalOptions

	// Layout behavior
	CenterContent bool
	FillHeight    bool
}

// ModalOptions defines modal overlay configuration
type ModalOptions struct {
	Visible bool
	Content string
	Style   *lipgloss.Style
}

// HelpContext provides context for help bar creation
type HelpContext struct {
	Mode          string
	ShowAll       bool
	CustomKeys    []string
	CustomHelp    map[string]string
	BubbleTeaHelp string // Pre-rendered Bubble Tea help content
}

// StatusContext provides context for status bar creation
type StatusContext struct {
	Mode          string
	Loading       bool
	Error         error
	StatusMessage string
	ContextInfo   map[string]string
	Counters      map[string]int
	FilterActive  bool
}

// LayoutState represents the current state of the layout
type LayoutState struct {
	// Dimensions
	TerminalWidth  int
	TerminalHeight int
	ContentWidth   int
	ContentHeight  int

	// Component states
	HelpBarVisible   bool
	StatusBarVisible bool
	ModalVisible     bool

	// Current context
	CurrentMode string
	CurrentPage string
}

// LayoutTemplate defines a reusable layout template
type LayoutTemplate struct {
	Name        string
	Description string
	Options     LayoutOptions
	ContentFunc func(content string, state LayoutState) string
}

// ResponsiveBreakpoint defines responsive layout breakpoints
type ResponsiveBreakpoint struct {
	MinWidth int
	MaxWidth int
	Config   Config
}
