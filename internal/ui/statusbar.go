package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Status bar component factory
func (m model) createStatusBar() string {
	w := lipgloss.Width

	// Context and mode indicator
	var modeIndicator string
	contextIcon := m.getContextIcon()
	contextName := m.getContextName()

	if m.showContextList {
		// Special indicator for context selection mode
		modeIndicator = modeAppsStyle.Render("🎯 CONTEXT SELECTION")
	} else {
		switch m.mode {
		case modeApps:
			modeIndicator = modeAppsStyle.Render(fmt.Sprintf("%s %s", contextIcon, contextName))
		case modeRevs:
			modeIndicator = modeRevisionsStyle.Render("REVISIONS")
		case modeContainers:
			modeIndicator = modeContainersStyle.Render("CONTAINERS")
		case modeEnvVars:
			modeIndicator = modeContainersStyle.Render("ENV VARS")
		case modeResourceGroups:
			modeIndicator = modeAppsStyle.Render("RESOURCE GROUPS")
		}
	}

	// Status indicator
	var statusIndicator string
	var filterActive bool

	// Check if any filter input is focused
	switch m.mode {
	case modeApps:
		filterActive = m.appsFilterInput.Focused()
	case modeRevs:
		filterActive = m.revisionsFilterInput.Focused()
	case modeContainers:
		filterActive = m.containersFilterInput.Focused()
	case modeEnvVars:
		filterActive = m.envVarsFilterInput.Focused()
	case modeResourceGroups:
		filterActive = m.resourceGroupsFilterInput.Focused()
	}

	if filterActive {
		statusIndicator = statusLoadingStyle.Render("Filtering")
	} else if m.loading {
		statusIndicator = statusLoadingStyle.Render("Loading " + m.spinner.View())
	} else if m.err != nil {
		statusIndicator = statusErrorStyle.Render("Error")
	} else {
		statusIndicator = statusReadyStyle.Render("Ready")
	}

	// Count indicator
	var countIndicator string
	if m.showContextList {
		// Show number of context options
		contextCount := len(m.createContextItems())
		if contextCount == 1 {
			countIndicator = countStyle.Render("1 Option")
		} else {
			countIndicator = countStyle.Render(fmt.Sprintf("%d Options", contextCount))
		}
	} else {
		switch m.mode {
		case modeApps:
			if len(m.apps) == 1 {
				countIndicator = countStyle.Render("1 App")
			} else {
				countIndicator = countStyle.Render(fmt.Sprintf("%d Apps", len(m.apps)))
			}
		case modeRevs:
			if len(m.revs) == 1 {
				countIndicator = countStyle.Render("1 Revision")
			} else {
				countIndicator = countStyle.Render(fmt.Sprintf("%d Revisions", len(m.revs)))
			}
		case modeContainers:
			if len(m.ctrs) == 1 {
				countIndicator = countStyle.Render("1 Container")
			} else {
				countIndicator = countStyle.Render(fmt.Sprintf("%d Containers", len(m.ctrs)))
			}
		case modeEnvVars:
			// Count environment variables for the current container
			envCount := 0
			if m.currentContainerName != "" {
				for _, ctr := range m.ctrs {
					if ctr.Name == m.currentContainerName {
						envCount = len(ctr.Env)
						break
					}
				}
			}
			if envCount == 1 {
				countIndicator = countStyle.Render("1 Env Var")
			} else {
				countIndicator = countStyle.Render(fmt.Sprintf("%d Env Vars", envCount))
			}
		case modeResourceGroups:
			if len(m.resourceGroups) == 1 {
				countIndicator = countStyle.Render("1 Resource Group")
			} else {
				countIndicator = countStyle.Render(fmt.Sprintf("%d Resource Groups", len(m.resourceGroups)))
			}
		}
	}

	// Context indicator (for deeper navigation levels)
	var contextIndicator string
	switch m.mode {
	case modeRevs:
		if appName := m.getCurrentAppName(); appName != "" {
			contextIndicator = contextStyle.Render("App: " + appName)
		}
	case modeContainers:
		if appName := m.getCurrentAppName(); appName != "" && m.currentRevName != "" {
			contextIndicator = contextStyle.Render(fmt.Sprintf("App: %s@%s", appName, m.currentRevName))
		}
	case modeEnvVars:
		if appName := m.getCurrentAppName(); appName != "" && m.currentRevName != "" && m.currentContainerName != "" {
			contextIndicator = contextStyle.Render(fmt.Sprintf("Container: %s@%s/%s", appName, m.currentRevName, m.currentContainerName))
		}
	}

	// Resource group indicator
	var rgIndicator string
	if m.rg != "" {
		rgIndicator = rgStyle.Render("RG: " + m.rg)
	}

	// Calculate widths for fixed elements
	fixedWidth := w(modeIndicator) + w(statusIndicator) + w(countIndicator) + w(rgIndicator)
	if contextIndicator != "" {
		fixedWidth += w(contextIndicator)
	}

	// Status message (expandable middle section)
	statusMessage := m.statusLine
	if statusMessage == "" {
		if m.err != nil {
			statusMessage = m.err.Error()
		} else if m.loading {
			switch m.mode {
			case modeApps:
				statusMessage = "Loading container apps..."
			case modeRevs:
				statusMessage = "Loading revisions..."
			case modeContainers:
				statusMessage = "Loading containers..."
			case modeEnvVars:
				statusMessage = "Loading environment variables..."
			case modeResourceGroups:
				statusMessage = "Loading resource groups..."
			}
		} else {
			// Show appropriate message based on current state
			if m.showContextList {
				statusMessage = "Use ↑/↓ to navigate, Enter to select, Esc to cancel"
			} else if !m.isAnyFilterActive() && !m.confirm.Visible {
				statusMessage = "Press : to switch context"
			} else {
				statusMessage = ""
			}
		}
	}

	// Create expandable status text
	statusVal := statusText.
		Width(max(0, m.termW-fixedWidth-4)). // Leave some margin
		Render(statusMessage)

	// Build the status bar
	var elements []string
	elements = append(elements, modeIndicator)
	elements = append(elements, statusIndicator)
	if contextIndicator != "" {
		elements = append(elements, contextIndicator)
	}
	elements = append(elements, statusVal)
	elements = append(elements, countIndicator)
	if rgIndicator != "" {
		elements = append(elements, rgIndicator)
	}

	bar := lipgloss.JoinHorizontal(lipgloss.Top, elements...)
	return statusBarStyle.Width(m.termW).Render(bar)
}

// Helper function for max
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getContextIcon returns the icon for the current context
func (m model) getContextIcon() string {
	switch m.context {
	case contextJobs:
		return "⚡"
	default:
		return "📦"
	}
}

// getContextName returns the display name for the current context
func (m model) getContextName() string {
	switch m.context {
	case contextJobs:
		return "JOBS"
	default:
		return "APPS"
	}
}
