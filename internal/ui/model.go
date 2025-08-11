package ui

import (
	"fmt"
	"os"

	models "az-tui/internal/models"

	tea "github.com/charmbracelet/bubbletea"
)

// ------------------------------ UI ------------------------------

type item models.ContainerApp

func (i item) Title() string       { return i.Name }
func (i item) Description() string { return i.ResourceGroup }
func (i item) FilterValue() string { return i.Name + " " + i.ResourceGroup }

type mode int

const (
	modeApps mode = iota
	modeRevs
	modeContainers
)

// ConfirmDialog holds state for a generic yes/no modal.
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
	json string

	// Simple selection state (just indices)
	selectedApp       int
	selectedRevision  int
	selectedContainer int

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

// messages for async commands

type noop struct{}

func InitialModel() model {
	return model{
		// Pure data
		apps: nil,
		revs: nil,
		ctrs: nil,
		json: "",

		// Simple selection state
		selectedApp:       0,
		selectedRevision:  0,
		selectedContainer: 0,

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
}

func appID(a models.ContainerApp) string { return fmt.Sprintf("%s/%s", a.ResourceGroup, a.Name) }

func revKey(appID, rev string) string { return appID + "@" + rev }

func (m model) currentApp() (models.ContainerApp, bool) {
	if m.selectedApp >= 0 && m.selectedApp < len(m.apps) {
		return m.apps[m.selectedApp], true
	}
	return models.ContainerApp{}, false
}

func (m *model) enterRevsFor(a models.ContainerApp) tea.Cmd {
	m.mode = modeRevs
	m.currentAppID = appID(a)

	// Reset revision selection
	m.selectedRevision = 0

	return LoadRevsCmd(a)
}

// Call when leaving revisions mode.
func (m *model) leaveRevs() {
	m.mode = modeApps
	m.currentAppID = ""
}

func (m model) withConfirm(text string, onYes func(model) (model, tea.Cmd), onNo func(model) (model, tea.Cmd)) model {
	m.confirm.Visible = true
	m.confirm.Text = text
	m.confirm.OnYes = onYes
	m.confirm.OnNo = onNo
	return m
}

func (m model) Init() tea.Cmd {
	return LoadAppsCmd(m.rg)
}
