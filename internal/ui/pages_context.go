package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Simple context item for single-line display
type simpleContextItem struct {
	id      string // Use string for flexibility with navigation contexts
	display string
	enabled bool
}

func (i simpleContextItem) FilterValue() string { return i.display }

// Custom delegate for single-line items
type contextDelegate struct{}

func (d contextDelegate) Height() int                             { return 1 }
func (d contextDelegate) Spacing() int                            { return 0 }
func (d contextDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d contextDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(simpleContextItem)
	if !ok {
		return
	}

	str := i.display
	if !i.enabled {
		str += " (Coming Soon)"
	}

	fn := lipgloss.NewStyle().PaddingLeft(4).Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("205")).
				Bold(true).
				Render("> " + strings.Join(s, " "))
		}
	}

	_, _ = fmt.Fprint(w, fn(str))
}

// createContextItems creates context items based on current mode
func (m model) createContextItems() []list.Item {
	switch m.mode {
	case modeResourceGroups:
		// From resource groups, can go to container apps or jobs (no resource group selected)
		return []list.Item{
			simpleContextItem{
				id:      string(contextApps),
				display: "📦 Container Apps",
				enabled: true,
			},
			simpleContextItem{
				id:      string(contextJobs),
				display: "⚡ Container App Jobs",
				enabled: false, // To be implemented
			},
		}

	case modeApps:
		// From apps, can switch between apps and jobs (preserve resource group selection)
		return []list.Item{
			simpleContextItem{
				id:      string(contextApps),
				display: "📦 Container Apps",
				enabled: true,
			},
			simpleContextItem{
				id:      string(contextJobs),
				display: "⚡ Container App Jobs",
				enabled: false, // To be implemented
			},
		}

	case modeRevs:
		// From revisions, can only go to revisions (preserve resource group and app selection)
		return []list.Item{
			simpleContextItem{
				id:      string(contextRevisions),
				display: "🔄 App Revisions",
				enabled: true,
			},
		}

	case modeContainers:
		// From containers, can only go to containers (preserve resource group, app, and revision selection)
		return []list.Item{
			simpleContextItem{
				id:      string(contextContainers),
				display: "🐳 Containers",
				enabled: true,
			},
		}

	case modeEnvVars:
		// From env vars, can only go to env vars (preserve all selections)
		return []list.Item{
			simpleContextItem{
				id:      string(contextEnvVars),
				display: "🔧 Environment Variables",
				enabled: true,
			},
		}

	default:
		// Fallback - show top-level contexts
		return []list.Item{
			simpleContextItem{
				id:      string(contextApps),
				display: "📦 Container Apps",
				enabled: true,
			},
			simpleContextItem{
				id:      string(contextJobs),
				display: "⚡ Container App Jobs",
				enabled: false,
			},
		}
	}
}

// createContextList creates and configures the context selection list
func (m model) createContextList() list.Model {
	items := m.createContextItems()

	// Use our custom single-line delegate
	delegate := contextDelegate{}

	// Calculate available height properly by measuring actual footer components
	helpBar := m.createHelpBar()
	statusBar := m.createStatusBar()
	helpBarHeight := lipgloss.Height(helpBar)
	statusBarHeight := lipgloss.Height(statusBar)

	listWidth := m.termW
	listHeight := m.termH - helpBarHeight - statusBarHeight

	contextListModel := list.New(items, delegate, listWidth, listHeight)
	contextListModel.Title = ""                 // No title
	contextListModel.SetShowStatusBar(false)    // We'll use the main status bar
	contextListModel.SetFilteringEnabled(false) // Disable filtering for simplicity
	contextListModel.SetShowHelp(false)         // Use our own help system
	contextListModel.SetShowTitle(false)        // Disable title completely
	contextListModel.SetShowPagination(false)   // Disable pagination

	return contextListModel
}

// handleContextListKey handles key events when the context list is visible
func (m model) handleContextListKey(msg tea.KeyMsg) (model, tea.Cmd) {
	switch msg.String() {
	case "?":
		// Handle help toggle in context mode
		m.help.ShowAll = !m.help.ShowAll
		// Recalculate context list size since help bar height changed
		m.contextList = m.createContextList()
		return m, nil

	case "enter":
		selectedItem, ok := m.contextList.SelectedItem().(simpleContextItem)
		if !ok {
			m.showContextList = false
			return m, nil
		}

		// Check if context is enabled
		if !selectedItem.enabled {
			m.statusLine = fmt.Sprintf("%s is not yet implemented", selectedItem.display)
			// Keep list open so user can select another option
			return m, nil
		}

		// Handle different types of selections
		switch selectedItem.id {
		case string(contextApps):
			// Navigate to container apps
			if m.mode == modeResourceGroups {
				// From resource groups to apps
				m.mode = modeApps
				m.context = contextApps
				m.loading = true
				m.apps = nil
				m.appsTable = m.createAppsTable()
				m.showContextList = false

				if m.rg != "" {
					// Load apps for specific resource group
					m.statusLine = "Loading Container Apps for " + m.rg
					return m, LoadAppsCmd(m.dataProvider, m.rg)
				} else {
					// Load all apps across all resource groups
					m.statusLine = "Loading All Container Apps"
					return m, LoadAppsCmd(m.dataProvider, "")
				}
			} else {
				// From apps mode - stay in apps (preserve resource group selection)
				m.mode = modeApps
				m.context = contextApps
				m.statusLine = "Container Apps"
			}

		case string(contextJobs):
			// Context switching to jobs (not implemented yet)
			m.statusLine = fmt.Sprintf("%s is not yet implemented", selectedItem.display)
			// Keep list open
			return m, nil

		case string(contextRevisions):
			// Stay in revisions mode (preserve resource group and app selection)
			m.mode = modeRevs
			m.statusLine = "App Revisions"

		case string(contextContainers):
			// Stay in containers mode (preserve resource group, app, and revision selection)
			m.mode = modeContainers
			m.statusLine = "Containers"

		case string(contextEnvVars):
			// Stay in env vars mode (preserve all selections)
			m.mode = modeEnvVars
			m.statusLine = "Environment Variables"
		}

		m.showContextList = false
		return m, nil

	case "esc", "q", ":":
		// Close without changing context
		m.showContextList = false
		return m, nil

	case "ctrl+c":
		// Allow quitting from context list
		return m, tea.Quit

	default:
		// Pass other keys to list component
		var cmd tea.Cmd
		m.contextList, cmd = m.contextList.Update(msg)
		return m, cmd
	}
}

// viewContextList renders the context selection list as a full-screen view
func (m model) viewContextList() string {
	// Just show the list view - no header, no current state
	listView := m.contextList.View()

	// Create the full layout with status bar
	layout := m.createTableLayout(listView)

	return layout
}
