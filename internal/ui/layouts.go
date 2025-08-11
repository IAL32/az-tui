package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// General layout manager for consistent UI structure
func (m model) createLayout(content string) string {
	// Create the main content area (content only, no title)
	mainContent := content

	// Create the bottom bars
	helpBar := m.createHelpBar()
	statusBar := m.createStatusBar()

	// Calculate available height for main content (total height - help bar - status bar)
	helpBarHeight := lipgloss.Height(helpBar)
	statusBarHeight := lipgloss.Height(statusBar)
	mainContentHeight := m.termH - helpBarHeight - statusBarHeight

	// Position main content at top, help bar and status bar at bottom
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Height(mainContentHeight).Render(mainContent),
		helpBar,
		statusBar,
	)

	if m.confirm.Visible {
		return lipgloss.Place(m.termW, m.termH, lipgloss.Center, lipgloss.Center, m.confirmBox())
	}
	return body
}

// Loading state layout
func (m model) createLoadingLayout(message string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styleAccent.Render(message),
		"",
	)
	return m.createLayout(content)
}

// Error state layout
func (m model) createErrorLayout(errorMsg string, helpMsg string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		StyleError.Render("Error: ")+errorMsg,
		styleAccent.Render(helpMsg),
		"",
	)
	return m.createLayout(content)
}

// Table layout
func (m model) createTableLayout(tableView string) string {
	return m.createLayout(tableView)
}
