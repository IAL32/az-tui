package ui

import (
	"fmt"
	"os"

	"github.com/IAL32/az-tui/internal/mock"
	models "github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/providers"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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

// Context types for switching between different views
type contextType string

const (
	contextApps contextType = "apps"
	contextJobs contextType = "jobs"
	// Navigation contexts
	contextResourceGroups contextType = "resource-groups"
	contextRevisions      contextType = "revisions"
	contextContainers     contextType = "containers"
	contextEnvVars        contextType = "env-vars"
)

// Types
type ConfirmDialog struct {
	Visible bool
	Text    string
	OnYes   func(m model) (model, tea.Cmd) // executed if user presses yes
	OnNo    func(m model) (model, tea.Cmd) // executed if user presses no/cancel
}

// Page Model Structures - Hierarchical State Management
type ResourceGroupsPageModel struct {
	Data        []models.ResourceGroup
	Table       table.Model
	FilterInput textinput.Model
	IsLoading   bool
	Error       error
}

type AppsPageModel struct {
	Data        []models.ContainerApp
	Table       table.Model
	FilterInput textinput.Model
	IsLoading   bool
	Error       error
}

type RevisionsPageModel struct {
	Data        []models.Revision
	Table       table.Model
	FilterInput textinput.Model
	IsLoading   bool
	Error       error
}

type ContainersPageModel struct {
	Data        []models.Container
	Table       table.Model
	FilterInput textinput.Model
	IsLoading   bool
	Error       error
}

type EnvVarsPageModel struct {
	Table       table.Model
	FilterInput textinput.Model
	IsLoading   bool
	Error       error
}

// Key bindings for different modes
type keyMap struct {
	// Navigation
	Enter     key.Binding
	Back      key.Binding
	Quit      key.Binding
	Up        key.Binding
	Down      key.Binding
	VimUp     key.Binding
	VimDown   key.Binding
	UpCombo   key.Binding // Combined up/k binding for help display
	DownCombo key.Binding // Combined down/j binding for help display

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

	// Context switching
	ContextSwitch key.Binding

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
		{k.ScrollLeft, k.ScrollRight, k.ContextSwitch, k.Help, k.Quit},
	}
}

// ContextHelp returns keybindings for context selection mode
func (k keyMap) ContextHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Enter, k.Back, k.Help, k.Quit},
	}
}

type model struct {
	// Global state - Context management
	context         contextType // Current active context
	contextList     list.Model  // List component for context selection
	showContextList bool        // Whether context list is visible

	// Page Models - Hierarchical State Management
	resourceGroupsPage ResourceGroupsPageModel
	appsPage           AppsPageModel
	revisionsPage      RevisionsPageModel
	containersPage     ContainersPageModel
	envVarsPage        EnvVarsPageModel

	// Global UI components
	help    help.Model
	keys    keyMap
	spinner spinner.Model

	// Global navigation context
	currentRG            string // Current resource group
	currentAppID         string // When viewing revisions
	currentRevName       string // When viewing containers
	currentContainerName string // When viewing environment variables

	// Performance cache (shared across pages)
	containersByRev map[string][]models.Container // key: revKey(appID, revName)

	// Global app state
	mode mode

	// Terminal dimensions (global)
	termW, termH int

	// Global status and confirmation
	statusLine string
	confirm    ConfirmDialog

	// Shared providers
	dataProvider    providers.DataProvider
	commandProvider providers.CommandProvider
}

// Messages

// Initialization
func InitialModel(useMockMode bool) model {
	// Initialize the appropriate data provider
	var dataProvider providers.DataProvider
	if useMockMode {
		mockProvider, err := mock.NewProvider()
		if err != nil {
			// Fallback to Azure provider if mock fails to load
			dataProvider = providers.NewAzureProvider()
		} else {
			dataProvider = mockProvider
		}
	} else {
		dataProvider = providers.NewAzureProvider()
	}

	m := model{
		// Global state - Context management
		context:         contextApps, // Default to apps
		showContextList: false,

		// Global navigation context
		currentRG:            os.Getenv("ACA_RG"),
		currentAppID:         "",
		currentRevName:       "",
		currentContainerName: "",

		// Performance cache (shared across pages)
		containersByRev: make(map[string][]models.Container),

		// Global app state
		mode: modeResourceGroups, // Start with resource groups mode

		// Terminal dimensions (global)
		termW: 80,
		termH: 24,

		// Global status and confirmation
		statusLine: "",
		confirm:    ConfirmDialog{},

		// Shared providers
		dataProvider:    dataProvider,
		commandProvider: createCommandProvider(useMockMode),
	}

	// Initialize context list using the dedicated method
	m.contextList = m.createContextList()

	// Initialize page models with their components
	m.resourceGroupsPage = ResourceGroupsPageModel{
		Data:        nil,
		Table:       m.createResourceGroupsTable(),
		FilterInput: createFilterInput("Filter resource groups..."),
		IsLoading:   true,
		Error:       nil,
	}

	m.appsPage = AppsPageModel{
		Data:        nil,
		Table:       m.createAppsTable(),
		FilterInput: createFilterInput("Filter apps..."),
		IsLoading:   false,
		Error:       nil,
	}

	m.revisionsPage = RevisionsPageModel{
		Data:        nil,
		Table:       m.createRevisionsTable(),
		FilterInput: createFilterInput("Filter revisions..."),
		IsLoading:   false,
		Error:       nil,
	}

	m.containersPage = ContainersPageModel{
		Data:        nil,
		Table:       m.createContainersTable(),
		FilterInput: createFilterInput("Filter containers..."),
		IsLoading:   false,
		Error:       nil,
	}

	m.envVarsPage = EnvVarsPageModel{
		Table:       m.createEnvVarsTable(),
		FilterInput: createFilterInput("Filter environment variables..."),
		IsLoading:   false,
		Error:       nil,
	}

	// Initialize global UI components
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
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		VimUp: key.NewBinding(
			key.WithKeys("k"),
			key.WithHelp("k", "up"),
		),
		VimDown: key.NewBinding(
			key.WithKeys("j"),
			key.WithHelp("j", "down"),
		),
		UpCombo: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		DownCombo: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
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
		ContextSwitch: key.NewBinding(
			key.WithKeys(":"),
			key.WithHelp(":", "switch context"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
	}

	return m
}

// Helper function to create filter input with placeholder
func createFilterInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	return input
}

func (m model) Init() tea.Cmd {
	return tea.Batch(LoadResourceGroupsCmd(m.dataProvider), m.spinner.Tick)
}

// createCommandProvider creates the appropriate command provider based on mock mode
func createCommandProvider(useMockMode bool) providers.CommandProvider {
	if useMockMode {
		return providers.NewMockCommandProvider()
	}
	return providers.NewAzureCommandProvider()
}

// Helper functions
func appID(a models.ContainerApp) string {
	return fmt.Sprintf("%s/%s", a.ResourceGroup, a.Name)
}

func revKey(appID, rev string) string {
	return appID + "@" + rev
}

func (m model) currentApp() (models.ContainerApp, bool) {
	if len(m.appsPage.Data) == 0 {
		return models.ContainerApp{}, false
	}

	// Get selected app from table
	selectedRow := m.appsPage.Table.HighlightedRow()
	if selectedRow.Data == nil {
		return models.ContainerApp{}, false
	}

	appName, ok := selectedRow.Data[columnKeyAppName].(string)
	if !ok {
		return models.ContainerApp{}, false
	}

	// Find the app by name
	for _, app := range m.appsPage.Data {
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
