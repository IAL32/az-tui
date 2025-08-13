package ui

import (
	"github.com/IAL32/az-tui/internal/ui/layouts"
	"github.com/charmbracelet/lipgloss"
)

// Legacy style variables for backward compatibility
// These are now managed through the layout system but kept for existing code
var (
	StyleError  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleAccent = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))
)

// GetLayoutManager returns a configured layout manager with default theme and dimensions
func GetLayoutManager(termW, termH int) *layouts.Manager {
	return layouts.NewManager(termW, termH, layouts.DefaultConfig())
}

// GetDefaultTheme returns the default theme configuration
func GetDefaultTheme() layouts.ThemeConfig {
	return layouts.DefaultTheme()
}

// Legacy style getters for backward compatibility
func GetStatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
		Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})
}

func GetStatusNuggetStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1)
}

func GetModeStyle(mode string) lipgloss.Style {
	base := GetStatusNuggetStyle().Bold(true)

	switch mode {
	case "apps":
		return base.Background(lipgloss.Color("#874BFD"))
	case "revisions":
		return base.Background(lipgloss.Color("#43BF6D"))
	case "containers":
		return base.Background(lipgloss.Color("#FF8C00"))
	case "resource-groups":
		return base.Background(lipgloss.Color("#20B2AA"))
	case "env-vars":
		return base.Background(lipgloss.Color("#9370DB"))
	default:
		return base.Background(lipgloss.Color("#888"))
	}
}

func GetStatusStyle(status string) lipgloss.Style {
	base := GetStatusNuggetStyle()

	switch status {
	case "loading":
		return base.Background(lipgloss.Color("#FFB347"))
	case "ready":
		return base.Background(lipgloss.Color("#32CD32"))
	case "error":
		return base.Background(lipgloss.Color("#FF6B6B"))
	default:
		return base.Background(lipgloss.Color("#6495ED"))
	}
}

func GetTableBaseStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a7a")).
		BorderForeground(lipgloss.Color("#a38"))
}

// GetThemeManager returns a theme manager with the default theme
func GetThemeManager() *layouts.ThemeManager {
	return layouts.NewThemeManager(layouts.DefaultTheme())
}
