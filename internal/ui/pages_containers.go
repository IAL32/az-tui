package ui

import (
	models "github.com/IAL32/az-tui/internal/models"

	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// createContainersTable creates a table component for containers
func (m model) createContainersTable() table.Model {
	// Create dynamic column builder
	builder := NewDynamicColumnBuilder().
		AddColumn(columnKeyCtrName, "Container", 12, true).       // Dynamic width, min 12
		AddColumn(columnKeyCtrStatus, "Status", 10, true).        // Fixed width - moved to second position
		AddColumn(columnKeyCtrImage, "Image", 50, true).          // Fixed width (longest content)
		AddColumn(columnKeyCtrCommand, "Command", 25, false).     // Fixed width
		AddColumn(columnKeyCtrArgs, "Args", 25, false).           // Fixed width
		AddColumn(columnKeyCtrResources, "Resources", 15, false). // Fixed width
		AddColumn(columnKeyCtrEnvCount, "Env", 8, false).         // Fixed width
		AddColumn(columnKeyCtrProbes, "Probes", 12, false).       // Fixed width
		AddColumn(columnKeyCtrVolumes, "Volumes", 10, false)      // Fixed width

	// Update dynamic column widths based on actual content
	for _, ctr := range m.ctrs {
		builder.UpdateWidthFromString(columnKeyCtrName, ctr.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(m.ctrs) > 0 {
		rows = make([]table.Row, len(m.ctrs))
		for i, ctr := range m.ctrs {
			command := strings.Join(ctr.Command, " ")
			if command == "" {
				command = "-"
			}

			args := strings.Join(ctr.Args, " ")
			if args == "" {
				args = "-"
			}

			// Resources
			resources := fmt.Sprintf("%.2gC/%.1s", ctr.CPU, ctr.Memory)
			if ctr.CPU == 0 {
				resources = "-"
			}

			// Environment variables count
			envCount := fmt.Sprintf("%d", len(ctr.Env))
			if len(ctr.Env) == 0 {
				envCount = "-"
			}

			// Probes
			probes := strings.Join(ctr.Probes, ",")
			if probes == "" {
				probes = "-"
			}

			// Volume mounts count
			volumes := fmt.Sprintf("%d", len(ctr.VolumeMounts))
			if len(ctr.VolumeMounts) == 0 {
				volumes = "-"
			}

			rows[i] = table.NewRow(table.RowData{
				columnKeyCtrName:      ctr.Name,
				columnKeyCtrImage:     ctr.Image,
				columnKeyCtrCommand:   command,
				columnKeyCtrArgs:      args,
				columnKeyCtrResources: resources,
				columnKeyCtrEnvCount:  envCount,
				columnKeyCtrProbes:    probes,
				columnKeyCtrVolumes:   volumes,
				columnKeyCtrStatus:    table.NewStyledCell("Running", lipgloss.NewStyle().Foreground(lipgloss.Color(getStatusColor("Running")))),
			})
		}
	} else {
		// Create empty table with placeholder
		rows = []table.Row{
			table.NewRow(table.RowData{
				columnKeyCtrName:      "No containers",
				columnKeyCtrImage:     "",
				columnKeyCtrCommand:   "",
				columnKeyCtrArgs:      "",
				columnKeyCtrResources: "",
				columnKeyCtrEnvCount:  "",
				columnKeyCtrProbes:    "",
				columnKeyCtrVolumes:   "",
				columnKeyCtrStatus:    "",
			}),
		}
	}

	t := table.New(columns).
		WithRows(rows).
		BorderRounded().
		WithBaseStyle(lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a7a")).
			BorderForeground(lipgloss.Color("#a38"))).
		WithMaxTotalWidth(m.termW).
		WithHorizontalFreezeColumnCount(1).
		Filtered(true).
		WithFilterInput(m.containersFilterInput).
		Focused(true)

	// Calculate height dynamically based on actual help and status bar heights
	helpBar := m.createHelpBar()
	statusBar := m.createStatusBar()
	helpBarHeight := lipgloss.Height(helpBar)
	statusBarHeight := lipgloss.Height(statusBar)

	// Available height = total height - help bar - status bar
	availableHeight := m.termH - helpBarHeight - statusBarHeight
	if availableHeight > 0 {
		t = t.WithPageSize(availableHeight)
	}

	return t
}

// getContainersHelpKeys returns the key bindings for containers mode
func (m model) getContainersHelpKeys() keyMap {
	return keyMap{
		Refresh:     m.keys.Refresh,
		Filter:      m.keys.Filter,
		Logs:        m.keys.Logs,
		Exec:        m.keys.Exec,
		EnvVars:     m.keys.EnvVars,
		ScrollLeft:  m.keys.ScrollLeft,
		ScrollRight: m.keys.ScrollRight,
		Help:        m.keys.Help,
		Back:        m.keys.Back,
		Quit:        m.keys.Quit,
	}
}

// Navigation functions
func (m *model) leaveContainers() {
	m.mode = modeRevs
	m.currentRevName = ""

	// Clear containers state
	m.ctrs = nil
	m.containersTable = m.createContainersTable()
}

// Message handlers
func (m model) handleLoadedContainersMsg(msg loadedContainersMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.ctrs = nil
		return m, nil
	}

	// cache
	m.containersByRev[revKey(msg.appID, msg.revName)] = msg.ctrs
	m.ctrs = msg.ctrs

	// Create the containers table with the loaded data
	if len(m.ctrs) > 0 {
		m.containersTable = m.createContainersTable()
	}

	return m, nil
}

// Key handlers
func (m model) handleContainersKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	// Handle filter input when focused
	if m.containersFilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.containersFilterInput.Blur()
			// Sync the filter with the table after applying
			m.containersTable = m.containersTable.WithFilterInput(m.containersFilterInput)
			return m, nil, true
		case "esc":
			m.containersFilterInput.SetValue("")
			m.containersFilterInput.Blur()
			m.containersTable = m.containersTable.WithFilterInput(m.containersFilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.containersFilterInput, cmd = m.containersFilterInput.Update(msg)
			m.containersTable = m.containersTable.WithFilterInput(m.containersFilterInput)
			return m, cmd, true
		}
	}

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit, true
	case "?":
		m.help.ShowAll = !m.help.ShowAll
		return m, nil, true
	case "/":
		m.containersFilterInput.SetValue("") // Clear any existing value
		m.containersFilterInput.Focus()
		m.containersTable = m.containersTable.WithFilterInput(m.containersFilterInput)
		return m, nil, true
	case "r":
		// Refresh containers list - clear data and show loading state
		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}
		m.ctrs = nil
		m.containersTable = m.createContainersTable()
		return m, LoadContainersCmd(a, m.currentRevName), true
	case "s":
		if len(m.ctrs) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		containerName, ok := selectedRow.Data[columnKeyCtrName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the container by name
		var selectedContainer models.Container
		found := false
		for _, ctr := range m.ctrs {
			if ctr.Name == containerName {
				selectedContainer = ctr
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}

		return m, m.azureCommands.ExecIntoContainer(a, m.currentRevName, selectedContainer.Name), true

	case "l":
		if len(m.ctrs) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		containerName, ok := selectedRow.Data[columnKeyCtrName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the container by name
		var selectedContainer models.Container
		found := false
		for _, ctr := range m.ctrs {
			if ctr.Name == containerName {
				selectedContainer = ctr
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}

		return m, m.azureCommands.ShowContainerLogs(a, m.currentRevName, selectedContainer.Name), true

	case "v":
		if len(m.ctrs) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersTable.HighlightedRow()
		if selectedRow.Data == nil {
			return m, nil, true
		}

		containerName, ok := selectedRow.Data[columnKeyCtrName].(string)
		if !ok {
			return m, nil, true
		}

		// Find the container by name
		var selectedContainer models.Container
		found := false
		for _, ctr := range m.ctrs {
			if ctr.Name == containerName {
				selectedContainer = ctr
				found = true
				break
			}
		}

		if !found {
			return m, nil, true
		}

		// Enter environment variables mode
		m.mode = modeEnvVars
		m.currentContainerName = selectedContainer.Name
		m.envVarsTable = m.createEnvVarsTable()
		return m, nil, true

	case "esc":
		m.leaveContainers()
		return m, nil, true
	}

	return m, nil, false
}

// View functions
func (m model) viewContainers() string {
	if len(m.ctrs) == 0 && m.err == nil {
		// Show loading state using generalized layout
		return m.createLoadingLayout("Loading containers...")
	}

	if m.err != nil && len(m.ctrs) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(m.err.Error(), "[esc] back  [q] quit")
	}

	// Show table view using generalized layout
	tableView := m.containersTable.View()
	return m.createTableLayout(tableView)
}
