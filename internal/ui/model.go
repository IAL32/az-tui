package ui

import (
	"fmt"
	"os"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// Constants
type mode int

const (
	modeApps mode = iota
	modeRevs
	modeContainers
)

// Types
type ConfirmDialog struct {
	Visible bool
	Text    string
	OnYes   func(m model) (model, tea.Cmd) // executed if user presses yes
	OnNo    func(m model) (model, tea.Cmd) // executed if user presses no/cancel
}

type model struct {
	// Pure data - no UI components
	apps []models.ContainerApp
	revs []models.Revision
	ctrs []models.Container

	// Table components for each mode
	appsTable       table.Model
	revisionsTable  table.Model
	containersTable table.Model

	// Filter text inputs for each mode
	appsFilterInput       textinput.Model
	revisionsFilterInput  textinput.Model
	containersFilterInput textinput.Model

	// Context for navigation
	currentAppID   string // When viewing revisions
	currentRevName string // When viewing containers

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
		apps: nil,
		revs: nil,
		ctrs: nil,

		// Context for navigation
		currentAppID:   "",
		currentRevName: "",

		// Caches
		containersByRev: make(map[string][]models.Container),

		// App state
		mode:    modeApps,
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

	// Initialize filter text inputs
	m.appsFilterInput = textinput.New()
	m.appsFilterInput.Placeholder = "Filter apps..."
	m.revisionsFilterInput = textinput.New()
	m.revisionsFilterInput.Placeholder = "Filter revisions..."
	m.containersFilterInput = textinput.New()
	m.containersFilterInput.Placeholder = "Filter containers..."

	return m
}

func (m model) Init() tea.Cmd {
	return LoadAppsCmd(m.rg)
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
