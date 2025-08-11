package ui

import (
	"bytes"
	"encoding/json"
	"fmt"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// App-related messages
type loadedAppsMsg struct {
	apps []models.ContainerApp
	err  error
}

type loadedDetailsMsg struct {
	json string
	err  error
}

// Handle loadedAppsMsg
func (m model) handleLoadedAppsMsg(msg loadedAppsMsg) (model, tea.Cmd) {
	m.loading = false
	m.err = msg.err
	if msg.err != nil {
		return m, nil
	}
	m.apps = msg.apps

	if len(m.apps) == 0 {
		return m, nil
	}

	// Set initial selection
	m.selectedApp = 0

	// Trigger initial load for selected app
	return m, tea.Batch(
		LoadDetailsCmd(m.apps[m.selectedApp]),
		LoadRevsCmd(m.apps[m.selectedApp]),
	)
}

// Handle loadedDetailsMsg
func (m model) handleLoadedDetailsMsg(msg loadedDetailsMsg) (model, tea.Cmd) {
	m.err = msg.err
	if msg.err != nil {
		m.json = ""
		return m, nil
	}

	// Pretty-print so indentation is stable
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(msg.json), "", "  "); err != nil {
		// fall back to raw if indent fails
		m.json = msg.json
	} else {
		m.json = buf.String()
	}

	return m, nil
}

// Handle key events when in Apps mode.
func (m model) handleAppsKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "j", "down":
		if m.selectedApp < len(m.apps)-1 {
			m.selectedApp++
			// Load details for newly selected app
			return m, tea.Batch(
				LoadDetailsCmd(m.apps[m.selectedApp]),
				LoadRevsCmd(m.apps[m.selectedApp]),
			), true
		}
		return m, nil, true
	case "k", "up":
		if m.selectedApp > 0 {
			m.selectedApp--
			// Load details for newly selected app
			return m, tea.Batch(
				LoadDetailsCmd(m.apps[m.selectedApp]),
				LoadRevsCmd(m.apps[m.selectedApp]),
			), true
		}
		return m, nil, true
	case "enter":
		if len(m.apps) == 0 {
			return m, nil, true
		}
		a := m.apps[m.selectedApp]
		return m, m.enterRevsFor(a), true
	case "R":
		if len(m.apps) == 0 {
			return m, nil, true
		}
		return m, LoadRevsCmd(m.apps[m.selectedApp]), true
	case "l":
		a := m.apps[m.selectedApp]
		return m, m.azureCommands.ShowAppLogs(a), true
	case "s", "e":
		a := m.apps[m.selectedApp]
		return m, m.azureCommands.ExecIntoApp(a), true
	}
	return m, nil, false
}

// createAppsList creates a list component for container apps
func (m model) createAppsList() list.Model {
	items := make([]list.Item, len(m.apps))
	for i, app := range m.apps {
		items[i] = item(app)
	}

	l := list.New(items, list.NewDefaultDelegate(), 32, m.termH-2)
	l.Title = "Container Apps"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	// Set selection based on data model
	if m.selectedApp >= 0 && m.selectedApp < len(items) {
		l.Select(m.selectedApp)
	}

	return l
}

func (m model) headerForCurrent() string {
	if len(m.apps) == 0 || m.selectedApp < 0 || m.selectedApp >= len(m.apps) {
		return ""
	}
	curr := m.apps[m.selectedApp]
	fqdn := curr.IngressFQDN
	if fqdn == "" {
		fqdn = "-"
	}
	return fmt.Sprintf("Name: %s  |  RG: %s  |  Loc: %s  |  FQDN: %s  |  Latest: %s",
		curr.Name, curr.ResourceGroup, curr.Location, fqdn, curr.LatestRevision)
}

func (m model) viewApps() string {
	if m.loading {
		spinner := m.createSpinner()
		return styleTitle.Render("Loading apps… ") + spinner.View()
	}
	if m.err != nil {
		return StyleError.Render("Error: ") + m.err.Error() + " Press r to retry or q to quit."
	}

	// Create components from current data
	appsList := m.createAppsList()
	detailsView := m.createDetailsView()
	revisionsTable := m.createRevisionsTable()

	left := appsList.View()
	right := lipgloss.JoinVertical(
		lipgloss.Left,
		styleTitle.Render("Details")+"\n"+detailsView.View(),
		styleTitle.Render("Revisions")+"\n"+revisionsTable.View(),
	)
	help := styleAccent.Render("[enter] revisions  [l] logs  [s] exec  [r] refresh  [R] reload revs  [q] quit  (j/k navigate)")

	body := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(34).Render(left),
		lipgloss.NewStyle().Padding(0, 1).Render(right),
	) + "\n" + help + "\n" + m.statusLine

	if m.confirm.Visible {
		return lipgloss.Place(m.termW, m.termH, lipgloss.Center, lipgloss.Center, m.confirmBox())
	}
	return body
}
