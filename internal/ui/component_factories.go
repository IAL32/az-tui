package ui

import (
	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

// Component factory functions - create UI components from data on-demand

// createDetailsView creates a viewport component for details
func (m model) createDetailsView() viewport.Model {
	rightW := max(20, m.termW-34-2)

	var height int
	if m.mode == modeApps {
		height = (m.termH - 4) / 2 // Top half of right pane
	} else {
		height = m.termH - 4 // Full right pane height
	}

	vp := viewport.New(rightW, height)
	vp.YPosition = 0

	content := m.getDetailsContent()
	vp.SetContent(content)

	return vp
}

// createSpinner creates a spinner component
func (m model) createSpinner() spinner.Model {
	sp := spinner.New()
	sp.Spinner = spinner.Dot
	return sp
}

// Helper methods for component factories

func (m model) getCurrentAppName() string {
	if len(m.apps) > 0 && m.selectedApp >= 0 && m.selectedApp < len(m.apps) {
		return m.apps[m.selectedApp].Name
	}
	return ""
}

func (m model) getCurrentApp() models.ContainerApp {
	if len(m.apps) > 0 && m.selectedApp >= 0 && m.selectedApp < len(m.apps) {
		return m.apps[m.selectedApp]
	}
	return models.ContainerApp{}
}

func (m model) getDetailsContent() string {
	switch m.mode {
	case modeApps:
		if len(m.apps) > 0 && m.selectedApp < len(m.apps) {
			if m.json != "" {
				return m.headerForCurrent() + "\n\n" + m.json
			}
			return m.headerForCurrent() + "\n\nLoading details..."
		}
		return "Select an app…"

	case modeRevs:
		return m.getRevisionContext()

	case modeContainers:
		if len(m.ctrs) > 0 && m.selectedContainer >= 0 && m.selectedContainer < len(m.ctrs) {
			app := m.getCurrentApp()
			header := m.containerHeader(app, m.currentRevName)
			content := header + "\n\n" + m.prettyContainerJSON(m.ctrs[m.selectedContainer])
			return content
		}
		app := m.getCurrentApp()
		if app.Name != "" {
			return m.containerHeader(app, m.currentRevName) + "\n\nLoading containers…"
		}
		return "No containers found."

	default:
		return ""
	}
}

// Helper methods for component factories

func (m model) confirmBox() string {
	if !m.confirm.Visible {
		return ""
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center).
		Render(m.confirm.Text + "\n\n[y] Yes  [n] No")

	return box
}
