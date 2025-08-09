package ui

import (
	"os"

	models "az-tui/internal/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// ------------------------------ UI ------------------------------

type pane int

const (
	paneDetails pane = iota
	paneRevisions
)

type item models.ContainerApp

func (i item) Title() string       { return i.Name }
func (i item) Description() string { return i.ResourceGroup }
func (i item) FilterValue() string { return i.Name + " " + i.ResourceGroup }

type model struct {
	// data
	apps []models.ContainerApp
	revs []models.Revision
	json string

	// selection
	cursor            int
	lastSelectedIndex int
	activePane        pane

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

	vp := viewport.New(80, 20)
	vp.YPosition = 0
	vp.SetContent("Select an app…")

	return model{
		apps:              nil,
		revs:              nil,
		json:              "",
		cursor:            0,
		lastSelectedIndex: -1,
		activePane:        paneDetails,
		rg:                os.Getenv("ACA_RG"),
		list:              l,
		spin:              sp,
		revTable:          t,
		jsonView:          vp,
		loading:           true,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(LoadAppsCmd(m.rg), m.spin.Tick)
}
