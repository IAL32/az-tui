package ui

import "github.com/charmbracelet/lipgloss"

var (
	styleTitle  = lipgloss.NewStyle().Bold(true)
	StyleError  = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	styleAccent = lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	// Status bar styles inspired by Lip Gloss example
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	// Mode indicators with different colors
	modeAppsStyle = statusNugget.
			Background(lipgloss.Color("#874BFD")).
			Bold(true)

	modeRevisionsStyle = statusNugget.
				Background(lipgloss.Color("#43BF6D")).
				Bold(true)

	modeContainersStyle = statusNugget.
				Background(lipgloss.Color("#FF8C00")).
				Bold(true)

	// Status indicators
	statusLoadingStyle = statusNugget.
				Background(lipgloss.Color("#FFB347"))

	statusReadyStyle = statusNugget.
				Background(lipgloss.Color("#32CD32"))

	statusErrorStyle = statusNugget.
				Background(lipgloss.Color("#FF6B6B"))

	// Info nuggets
	countStyle = statusNugget.
			Background(lipgloss.Color("#6495ED"))

	contextStyle = statusNugget.
			Background(lipgloss.Color("#9370DB"))

	rgStyle = statusNugget.
		Background(lipgloss.Color("#20B2AA"))

	// Status text (expandable middle section)
	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)
)
