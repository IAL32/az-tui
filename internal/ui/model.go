package ui

import (
	"fmt"
	"os"

	models "github.com/IAL32/az-tui/internal/models"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// Constants
type mode int

const (
	modeApps mode = iota
	modeRevs
	modeContainers
	modeEnvVars
	modeResourceGroups
)

// Types
type ConfirmDialog struct {
	Visible bool
	Text    string
	OnYes   func(m model) (model, tea.Cmd) // executed if user presses yes
	OnNo    func(m model) (model, tea.Cmd) // executed if user presses no/cancel
}

// Key bindings for different modes
type keyMap struct {
	// Navigation
	Enter key.Binding
	Back  key.Binding
	Quit  key.Binding

	// Actions
	Refresh    key.Binding
	RestartRev key.Binding
	Filter     key.Binding
	Logs       key.Binding
	Exec       key.Binding
	EnvVars    key.Binding

	// Table navigation
	ScrollLeft  key.Binding
	ScrollRight key.Binding

	// Help
	Help key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Refresh, k.Filter, k.Logs},
		{k.Exec, k.EnvVars, k.RestartRev, k.Back},
		{k.ScrollLeft, k.ScrollRight, k.Help, k.Quit},
	}
}

type model struct {
	// Pure data - no UI components
	apps           []models.ContainerApp
	revs           []models.Revision
	ctrs           []models.Container
	resourceGroups []models.ResourceGroup

	// Table components for each mode
	appsTable           table.Model
	revisionsTable      table.Model
	containersTable     table.Model
	envVarsTable        table.Model
	resourceGroupsTable table.Model

	// Filter text inputs for each mode
	appsFilterInput           textinput.Model
	revisionsFilterInput      textinput.Model
	containersFilterInput     textinput.Model
	envVarsFilterInput        textinput.Model
	resourceGroupsFilterInput textinput.Model

	// Help system
	help help.Model
	keys keyMap

	// Loading spinner
	spinner spinner.Model

	// Context for navigation
	currentAppID         string // When viewing revisions
	currentRevName       string // When viewing containers
	currentContainerName string // When viewing environment variables

	// Optional caches for performance
	containersByRev map[string][]models.Container // key: revKey(appID, revName)

	// App state
	mode    mode
	loading bool
	err     error

	// Configuration
	rg string

	// Terminal size for component sizing
	termW, termH int

	// Status and confirmation
	statusLine string
	confirm    ConfirmDialog

	// Command execution
	azureCommands *AzureCommands
}

// Messages
type noop struct{}

// Initialization
func InitialModel() model {
	m := model{
		// Pure data
		apps:           nil,
		revs:           nil,
		ctrs:           nil,
		resourceGroups: nil,

		// Context for navigation
		currentAppID:   "",
		currentRevName: "",

		// Caches
		containersByRev: make(map[string][]models.Container),

		// App state
		mode:    modeResourceGroups, // Start with resource groups mode
		loading: true,
		err:     nil,

		// Configuration
		rg: os.Getenv("ACA_RG"),

		// Terminal size (will be set on first WindowSizeMsg)
		termW: 80,
		termH: 24,

		// Status
		statusLine: "",
		confirm:    ConfirmDialog{},

		// Command execution
		azureCommands: NewAzureCommands(),
	}

	// Initialize empty tables
	m.appsTable = m.createAppsTable()
	m.revisionsTable = m.createRevisionsTable()
	m.containersTable = m.createContainersTable()
	m.envVarsTable = m.createEnvVarsTable()
	m.resourceGroupsTable = m.createResourceGroupsTable()

	// Initialize filter text inputs
	m.appsFilterInput = textinput.New()
	m.appsFilterInput.Placeholder = "Filter apps..."
	m.revisionsFilterInput = textinput.New()
	m.revisionsFilterInput.Placeholder = "Filter revisions..."
	m.containersFilterInput = textinput.New()
	m.containersFilterInput.Placeholder = "Filter containers..."
	m.envVarsFilterInput = textinput.New()
	m.envVarsFilterInput.Placeholder = "Filter environment variables..."
	m.resourceGroupsFilterInput = textinput.New()
	m.resourceGroupsFilterInput.Placeholder = "Filter resource groups..."

	// Initialize help system
	m.help = help.New()

	// Initialize spinner
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	m.keys = keyMap{
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
		RestartRev: key.NewBinding(
			key.WithKeys("R"),
			key.WithHelp("R", "restart revision"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		Logs: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "logs"),
		),
		Exec: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "exec"),
		),
		EnvVars: key.NewBinding(
			key.WithKeys("v"),
			key.WithHelp("v", "env vars"),
		),
		ScrollLeft: key.NewBinding(
			key.WithKeys("shift+left"),
			key.WithHelp("shift+←", "scroll left"),
		),
		ScrollRight: key.NewBinding(
			key.WithKeys("shift+right"),
			key.WithHelp("shift+→", "scroll right"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}

	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(LoadResourceGroupsCmd(), m.spinner.Tick)
}

// Helper functions
func appID(a models.ContainerApp) string {
	return fmt.Sprintf("%s/%s", a.ResourceGroup, a.Name)
}

func revKey(appID, rev string) string {
	return appID + "@" + rev
}

func (m model) currentApp() (models.ContainerApp, bool) {
	if len(m.apps) == 0 {
		return models.ContainerApp{}, false
	}

	// Get selected app from table
	selectedRow := m.appsTable.HighlightedRow()
	if selectedRow.Data == nil {
		return models.ContainerApp{}, false
	}

	appName, ok := selectedRow.Data[columnKeyAppName].(string)
	if !ok {
		return models.ContainerApp{}, false
	}

	// Find the app by name
	for _, app := range m.apps {
		if app.Name == appName {
			return app, true
		}
	}

	return models.ContainerApp{}, false
}

func (m model) withConfirm(text string, onYes func(model) (model, tea.Cmd), onNo func(model) (model, tea.Cmd)) model {
	m.confirm.Visible = true
	m.confirm.Text = text
	m.confirm.OnYes = onYes
	m.confirm.OnNo = onNo
	return m
}
