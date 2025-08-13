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
	for _, ctr := range m.containersPage.Data {
		builder.UpdateWidthFromString(columnKeyCtrName, ctr.Name)
	}

	// Build columns with calculated widths
	columns := builder.Build()

	var rows []table.Row
	if len(m.containersPage.Data) > 0 {
		rows = make([]table.Row, len(m.containersPage.Data))
		for i, ctr := range m.containersPage.Data {
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

	return m.createUnifiedTable(columns, rows, m.containersPage.FilterInput)
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
	m.containersPage.Data = nil
	m.containersPage.Table = m.createContainersTable()
}

// Message handlers
func (m model) handleLoadedContainersMsg(msg loadedContainersMsg) (model, tea.Cmd) {
	m.containersPage.IsLoading = false
	m.containersPage.Error = msg.err

	if msg.err != nil {
		m.containersPage.Data = nil
		return m, nil
	}

	// cache
	m.containersByRev[revKey(msg.appID, msg.revName)] = msg.ctrs
	m.containersPage.Data = msg.ctrs

	// Create the containers table with the loaded data
	m.containersPage.Table = m.createContainersTable()

	return m, nil
}

// Key handlers
func (m model) handleContainersKey(msg tea.KeyMsg) (model, tea.Cmd, bool) {
	// Handle filter input when focused
	if m.containersPage.FilterInput.Focused() {
		switch msg.String() {
		case "enter":
			m.containersPage.FilterInput.Blur()
			// Sync the filter with the table after applying
			m.containersPage.Table = m.containersPage.Table.WithFilterInput(m.containersPage.FilterInput)
			return m, nil, true
		case "esc":
			m.containersPage.FilterInput.SetValue("")
			m.containersPage.FilterInput.Blur()
			m.containersPage.Table = m.containersPage.Table.WithFilterInput(m.containersPage.FilterInput)
			return m, nil, true
		default:
			var cmd tea.Cmd
			m.containersPage.FilterInput, cmd = m.containersPage.FilterInput.Update(msg)
			m.containersPage.Table = m.containersPage.Table.WithFilterInput(m.containersPage.FilterInput)
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
		m.containersPage.FilterInput.SetValue("") // Clear any existing value
		m.containersPage.FilterInput.Focus()
		m.containersPage.Table = m.containersPage.Table.WithFilterInput(m.containersPage.FilterInput)
		return m, nil, true
	case "r":
		// Refresh containers list - clear data and show loading state
		a := m.getCurrentApp()
		if a.Name == "" || m.currentRevName == "" {
			return m, nil, true
		}
		m.containersPage.IsLoading = true
		m.containersPage.Error = nil
		m.containersPage.Data = nil
		m.containersPage.Table = m.createContainersTable()
		return m, LoadContainersCmd(m.dataProvider, a, m.currentRevName), true
	case "s":
		if len(m.containersPage.Data) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersPage.Table.HighlightedRow()
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
		for _, ctr := range m.containersPage.Data {
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

		return m, m.commandProvider.ExecIntoContainer(a, m.currentRevName, selectedContainer.Name), true

	case "l":
		if len(m.containersPage.Data) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersPage.Table.HighlightedRow()
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
		for _, ctr := range m.containersPage.Data {
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

		return m, m.commandProvider.ShowContainerLogs(a, m.currentRevName, selectedContainer.Name), true

	case "v":
		if len(m.containersPage.Data) == 0 {
			return m, nil, true
		}

		// Get selected container from table
		selectedRow := m.containersPage.Table.HighlightedRow()
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
		for _, ctr := range m.containersPage.Data {
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
		m.envVarsPage.Table = m.createEnvVarsTable()
		return m, nil, true

	case "esc":
		m.leaveContainers()
		return m, nil, true
	}

	return m, nil, false
}

// View functions
func (m model) viewContainers() string {
	if len(m.containersPage.Data) == 0 && m.containersPage.Error == nil {
		// Show loading state using generalized layout
		return m.createLoadingLayout("Loading containers...")
	}

	if m.containersPage.Error != nil && len(m.containersPage.Data) == 0 {
		// Show error state using generalized layout
		return m.createErrorLayout(m.containersPage.Error.Error(), "[esc] back  [q] quit")
	}

	// Show table view using generalized layout
	tableView := m.containersPage.Table.View()
	return m.createTableLayout(tableView)
}
