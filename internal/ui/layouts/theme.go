package layouts

import (
	"github.com/charmbracelet/lipgloss"
)

// ThemeConfig defines the theme configuration for the layout system
type ThemeConfig struct {
	// Color scheme
	Colors ColorScheme

	// Typography
	Typography TypographyConfig

	// Spacing
	Spacing SpacingConfig

	// Border styles
	Borders BorderConfig
}

// ColorScheme defines the color palette for the theme
type ColorScheme struct {
	// Primary colors
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color

	// Status colors
	Success lipgloss.Color
	Warning lipgloss.Color
	Error   lipgloss.Color
	Info    lipgloss.Color

	// UI colors
	Background lipgloss.Color
	Foreground lipgloss.Color
	Border     lipgloss.Color
	Highlight  lipgloss.Color

	// Adaptive colors for light/dark mode
	AdaptiveBackground lipgloss.AdaptiveColor
	AdaptiveForeground lipgloss.AdaptiveColor
	AdaptiveBorder     lipgloss.AdaptiveColor
}

// TypographyConfig defines typography settings
type TypographyConfig struct {
	FontFamily string
	FontSize   int
	LineHeight float64
	Bold       bool
	Italic     bool
}

// SpacingConfig defines spacing constants
type SpacingConfig struct {
	Margin  int
	Padding int
	Gap     int
}

// BorderConfig defines border styles
type BorderConfig struct {
	Style lipgloss.Border
	Width int
	Color lipgloss.Color
}

// ThemeManager manages theme configuration and provides styled components
type ThemeManager struct {
	config ThemeConfig
	styles map[string]lipgloss.Style
}

// NewThemeManager creates a new theme manager with the given configuration
func NewThemeManager(config ThemeConfig) *ThemeManager {
	tm := &ThemeManager{
		config: config,
		styles: make(map[string]lipgloss.Style),
	}
	tm.initializeStyles()
	return tm
}

// initializeStyles creates the default styles based on the theme configuration
func (tm *ThemeManager) initializeStyles() {
	// Status bar styles
	tm.styles["statusBar"] = lipgloss.NewStyle().
		Foreground(tm.config.Colors.AdaptiveForeground).
		Background(tm.config.Colors.AdaptiveBackground)

	tm.styles["statusNugget"] = lipgloss.NewStyle().
		Foreground(tm.config.Colors.Foreground).
		Padding(0, 1)

	// Mode indicators
	tm.styles["modeApps"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Primary).
		Bold(true)

	tm.styles["modeRevisions"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Success).
		Bold(true)

	tm.styles["modeContainers"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Warning).
		Bold(true)

	// Status indicators
	tm.styles["statusLoading"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Warning)

	tm.styles["statusReady"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Success)

	tm.styles["statusError"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Error)

	// Info nuggets
	tm.styles["count"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Info)

	tm.styles["context"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Secondary)

	tm.styles["resourceGroup"] = tm.styles["statusNugget"].Copy().
		Background(tm.config.Colors.Accent)

	// Status text
	tm.styles["statusText"] = tm.styles["statusBar"].Copy()

	// Table styles
	tm.styles["tableBase"] = lipgloss.NewStyle().
		Foreground(tm.config.Colors.Foreground).
		BorderForeground(tm.config.Colors.Border)

	// Error and accent styles
	tm.styles["error"] = lipgloss.NewStyle().
		Foreground(tm.config.Colors.Error)

	tm.styles["accent"] = lipgloss.NewStyle().
		Foreground(tm.config.Colors.Accent)
}

// GetStyle returns a style by name
func (tm *ThemeManager) GetStyle(name string) lipgloss.Style {
	if style, exists := tm.styles[name]; exists {
		return style
	}
	return lipgloss.NewStyle()
}

// UpdateStyle updates or creates a style
func (tm *ThemeManager) UpdateStyle(name string, style lipgloss.Style) {
	tm.styles[name] = style
}

// GetTheme returns the current theme configuration
func (tm *ThemeManager) GetTheme() ThemeConfig {
	return tm.config
}

// SetTheme updates the theme configuration and reinitializes styles
func (tm *ThemeManager) SetTheme(config ThemeConfig) {
	tm.config = config
	tm.initializeStyles()
}

// DefaultTheme returns the default theme configuration
func DefaultTheme() ThemeConfig {
	return ThemeConfig{
		Colors: ColorScheme{
			Primary:    lipgloss.Color("#874BFD"),
			Secondary:  lipgloss.Color("#9370DB"),
			Accent:     lipgloss.Color("#212"),
			Success:    lipgloss.Color("#32CD32"),
			Warning:    lipgloss.Color("#FFB347"),
			Error:      lipgloss.Color("#FF6B6B"),
			Info:       lipgloss.Color("#6495ED"),
			Background: lipgloss.Color("#FFFDF5"),
			Foreground: lipgloss.Color("#888888"),
			Border:     lipgloss.Color("#a38"),
			Highlight:  lipgloss.Color("#888888"),
			AdaptiveBackground: lipgloss.AdaptiveColor{
				Light: "#D9DCCF",
				Dark:  "#353533",
			},
			AdaptiveForeground: lipgloss.AdaptiveColor{
				Light: "#888888",
				Dark:  "#C1C6B2",
			},
			AdaptiveBorder: lipgloss.AdaptiveColor{
				Light: "#a38",
				Dark:  "#a7a",
			},
		},
		Typography: TypographyConfig{
			FontFamily: "monospace",
			FontSize:   12,
			LineHeight: 1.2,
			Bold:       false,
			Italic:     false,
		},
		Spacing: SpacingConfig{
			Margin:  1,
			Padding: 1,
			Gap:     1,
		},
		Borders: BorderConfig{
			Style: lipgloss.NormalBorder(),
			Width: 1,
			Color: lipgloss.Color("#a38"),
		},
	}
}
