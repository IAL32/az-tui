package ui

import (
	"fmt"
	"os"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ------------------------------ UI ------------------------------

type item models.ContainerApp

func (i item) Title() string       { return i.Name }
func (i item) Description() string { return i.ResourceGroup }
func (i item) FilterValue() string { return i.Name + " " + i.ResourceGroup }

type pane int
type mode int

const (
	paneDetails pane = iota
	paneRevisions
)
const (
	modeApps mode = iota
	modeRevs
	modeContainers
)

type model struct {
	// data
	apps []models.ContainerApp
	revs []models.Revision
	ctrs []models.Container
	json string

	// selection
	// Independent cursors + last-selected tracking

	appsCursor, lastAppsIndex int
	revsCursor, lastRevsIndex int
	ctrCursor, lastCtrIndex   int

	// Per-app revision cursor memory (restore when you return)
	revCursorByAppID map[string]int                // "rg/name" -> rev index    // Optional caches
	containersByRev  map[string][]models.Container // key: revKey(appID, revName)

	activePane pane

	// deps/config
	rg string

	// ui components
	list     list.Model
	spin     spinner.Model
	revTable table.Model
	jsonView viewport.Model

	// status
	loading bool
	err     error

	// mode
	mode     mode
	revList  list.Model
	ctrList  list.Model
	revAppID string
	revName  string // selected revision name (containers page context)
}

// messages for async commands

type loadedAppsMsg struct {
	apps []models.ContainerApp
	err  error
}
type loadedDetailsMsg struct {
	json string
	err  error
}
type loadedRevsMsg struct {
	revs []models.Revision
	err  error
}

type loadedContainersMsg struct {
	appID   string
	revName string
	ctrs    []models.Container
	err     error
}

type noop struct{}

func InitialModel() model {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 32, 20)
	l.Title = "Container Apps"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	sp := spinner.New()
	sp.Spinner = spinner.Dot

	t := table.New(
		table.WithColumns(
			[]table.Column{
				{Title: "Revision", Width: 28},
				{Title: "Active", Width: 6},
				{Title: "Traffic", Width: 8},
				{Title: "Created", Width: 18},
				{Title: "Status", Width: 12},
			},
		),
	)
	revList := list.New([]list.Item{}, list.NewDefaultDelegate(), 32, 20)
	revList.Title = "Revisions"
	revList.SetShowTitle(true)
	revList.SetShowStatusBar(false)
	revList.SetFilteringEnabled(true)

	ctrList := list.New([]list.Item{}, list.NewDefaultDelegate(), 32, 20)
	ctrList.Title = "Containers"
	ctrList.SetShowStatusBar(false)
	ctrList.SetFilteringEnabled(true)

	vp := viewport.New(80, 20)
	vp.YPosition = 0
	vp.SetContent("Select an app…")

	return model{
		apps:             nil,
		revs:             nil,
		ctrs:             nil,
		json:             "",
		appsCursor:       0,
		lastAppsIndex:    -1,
		revsCursor:       0,
		lastRevsIndex:    -1,
		ctrCursor:        0,
		lastCtrIndex:     -1,
		revCursorByAppID: make(map[string]int),
		containersByRev:  make(map[string][]models.Container),
		activePane:       paneDetails,
		rg:               os.Getenv("ACA_RG"),
		list:             l,
		spin:             sp,
		revTable:         t,
		jsonView:         vp,
		loading:          true,
		mode:             modeApps,
		revList:          revList,
		ctrList:          ctrList,
		revAppID:         "",
		revName:          "",
	}
}

func appID(a models.ContainerApp) string { return fmt.Sprintf("%s/%s", a.ResourceGroup, a.Name) }

func revKey(appID, rev string) string { return appID + "@" + rev }

func (m model) currentApp() (models.ContainerApp, bool) {
	if m.appsCursor >= 0 && m.appsCursor < len(m.apps) {
		return m.apps[m.appsCursor], true
	}
	return models.ContainerApp{}, false
}

func (m *model) syncAppsCursorFromList() {
	idx := m.list.Index()
	if idx >= 0 && idx < len(m.apps) {
		m.appsCursor = idx
	}
}

func (m *model) syncRevsCursorFromList() {
	idx := m.revList.Index()
	if idx >= 0 && idx < len(m.revs) {
		m.revsCursor = idx
	}
}

func (m *model) enterRevsFor(a models.ContainerApp) tea.Cmd {
	m.mode = modeRevs
	m.revAppID = appID(a)

	// Title reflects context
	m.revList.Title = fmt.Sprintf("Revisions — %s", a.Name)
	m.revList.SetShowTitle(true)

	// Mirror the current left pane size so the title has space immediately
	m.revList.SetSize(m.list.Width(), m.list.Height())

	m.revList.SetItems(nil)
	m.revTable.SetRows(nil) // right bottom pane
	return LoadRevsCmd(a)
}

// Call when leaving revisions mode.
func (m *model) leaveRevs() {
	// Remember current rev cursor for this app so we can restore next time.
	if m.revAppID != "" {
		m.revCursorByAppID[m.revAppID] = m.revsCursor
	}
	m.mode = modeApps
	m.revAppID = ""
	m.revList.SetItems(nil)
}

func (m model) Init() tea.Cmd {
	return tea.Batch(LoadAppsCmd(m.rg), m.spin.Tick)
}
