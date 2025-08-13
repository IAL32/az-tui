package ui

import (
	"fmt"
	"io"
	"strings"

	"github.com/IAL32/az-tui/internal/mock"
	"github.com/IAL32/az-tui/internal/providers"
	"github.com/IAL32/az-tui/internal/ui/core"
	"github.com/IAL32/az-tui/internal/ui/layouts"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
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

// ConfirmDialog represents a confirmation dialog
type ConfirmDialog struct {
	Visible bool
	Text    string
	OnYes   func(m model) (model, tea.Cmd) // executed if user presses yes
	OnNo    func(m model) (model, tea.Cmd) // executed if user presses no/cancel
}

// Key bindings for the main model
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

// model represents the main UI model that delegates to the core system
type model struct {
	// Core coordination system
	core core.CoreInterface

	// Global UI components
	help    help.Model
	keys    keyMap
	spinner spinner.Model

	// Terminal dimensions
	termW, termH int

	// Global status and confirmation
	confirm ConfirmDialog
}

// InitialModel creates the initial model with core coordination
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

	// Create command provider
	commandProvider := createCommandProvider(useMockMode)

	// Initialize terminal dimensions
	termW, termH := 80, 24

	// Create core model
	coreModel := core.NewCoreModel(dataProvider, commandProvider, termW, termH)

	// Create main model
	m := model{
		core:    coreModel,
		termW:   termW,
		termH:   termH,
		confirm: ConfirmDialog{},
	}

	// Initialize global UI components
	m.help = help.New()

	// Style the help component with theme colors
	theme := layouts.DefaultTheme()
	m.help.Styles.FullKey = lipgloss.NewStyle().Foreground(theme.Colors.Primary).Bold(true)
	m.help.Styles.FullDesc = lipgloss.NewStyle().Foreground(theme.Colors.AdaptiveForeground)
	m.help.Styles.FullSeparator = lipgloss.NewStyle().Foreground(theme.Colors.Border)
	m.help.Styles.ShortKey = lipgloss.NewStyle().Foreground(theme.Colors.Primary).Bold(true)
	m.help.Styles.ShortDesc = lipgloss.NewStyle().Foreground(theme.Colors.AdaptiveForeground)
	m.help.Styles.ShortSeparator = lipgloss.NewStyle().Foreground(theme.Colors.Border)

	// Initialize spinner
	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Dot
	m.spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Initialize key bindings
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

	// Initialize context list
	m.core.SetContextList(m.createContextList())

	return m
}

// Init initializes the model
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.core.LoadResourceGroups(),
		m.spinner.Tick,
	)
}

// Update handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m.handleWindowSize(msg)
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	default:
		// Handle spinner updates and delegate to core
		var spinnerCmd tea.Cmd
		m.spinner, spinnerCmd = m.spinner.Update(msg)
		coreCmd := m.core.HandleMessage(msg)

		// Batch non-nil commands efficiently
		var cmds []tea.Cmd
		if spinnerCmd != nil {
			cmds = append(cmds, spinnerCmd)
		}
		if coreCmd != nil {
			cmds = append(cmds, coreCmd)
		}

		if len(cmds) > 0 {
			return m, tea.Batch(cmds...)
		}
		return m, nil
	}
}

// View renders the model
func (m model) View() string {
	// Create help context from main model state
	helpContext := layouts.HelpContext{
		ShowAll: m.help.ShowAll,
	}

	// If help is showing, pre-render the Bubble Tea help and pass it to the layout
	if m.help.ShowAll {
		keyMap := m.getContextKeyMap()
		helpContext.BubbleTeaHelp = m.help.View(keyMap)
	}

	// Show context list when active
	if m.core.IsShowingContextList() {
		return m.viewContextList()
	}

	// Pass help context to core view (layout will handle positioning)
	return m.core.ViewWithHelpContext(helpContext)
}

// Message handlers

func (m model) handleWindowSize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	w, h := msg.Width, msg.Height
	m.termW, m.termH = w, h
	m.help.Width = w

	// Update core dimensions
	m.core.UpdateDimensions(w, h)

	// Update context list size
	listWidth := min(70, w-10)
	listHeight := min(12, h-6)
	contextList := m.core.GetContextList()
	contextList.SetSize(listWidth, listHeight)
	m.core.SetContextList(contextList)

	return m, nil
}

func (m model) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle confirmation dialog
	if m.confirm.Visible {
		switch msg.String() {
		case "y", "enter":
			m.confirm.Visible = false
			if m.confirm.OnYes != nil {
				return m.confirm.OnYes(m)
			}
			return m, nil
		case "n", "esc":
			m.confirm.Visible = false
			if m.confirm.OnNo != nil {
				return m.confirm.OnNo(m)
			}
			return m, nil
		}
		return m, nil // swallow all other keys when modal visible
	}

	// Handle context list when visible
	if m.core.IsShowingContextList() {
		return m.handleContextListKey(msg)
	}

	// Intercept ":" to show context list
	if msg.String() == ":" {
		// Don't show if in filter mode or confirm dialog
		if m.core.IsAnyFilterActive() || m.confirm.Visible {
			return m, nil
		}
		m.core.SetShowContextList(true)
		m.core.SetStatusLine("") // Clear status line when entering context mode
		// Recreate context list based on current mode
		m.core.SetContextList(m.createContextList())
		return m, nil
	}

	// Handle help toggle at the main model level
	if msg.String() == "?" {
		m.help.ShowAll = !m.help.ShowAll
		return m, nil
	}

	// Delegate key handling to core
	var cmd tea.Cmd
	var handled bool

	cmd, handled = m.core.HandleKeyMsg(msg)
	if handled {
		return m, cmd
	}

	// Handle table navigation for unhandled keys
	cmd = m.core.UpdateTable(msg)
	if cmd != nil {
		return m, cmd
	}

	return m, nil
}

// Context list management

func (m model) createContextList() list.Model {
	items := m.createContextItems()

	// Use our custom single-line delegate
	delegate := contextDelegate{}

	// Calculate available height properly by estimating footer components
	// Use standard heights for help and status bars
	helpBarHeight := 2
	statusBarHeight := 1

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

func (m model) createContextItems() []list.Item {
	switch m.core.GetCurrentMode() {
	case core.ModeResourceGroups:
		// From resource groups, can go to container apps or jobs (no resource group selected)
		return []list.Item{
			simpleContextItem{
				id:      "apps",
				display: "📦 Container Apps",
				enabled: true,
			},
			simpleContextItem{
				id:      "jobs",
				display: "⚡ Container App Jobs",
				enabled: false, // To be implemented
			},
		}

	case core.ModeApps:
		// From apps, can switch between apps and jobs (preserve resource group selection)
		return []list.Item{
			simpleContextItem{
				id:      "apps",
				display: "📦 Container Apps",
				enabled: true,
			},
			simpleContextItem{
				id:      "jobs",
				display: "⚡ Container App Jobs",
				enabled: false, // To be implemented
			},
		}

	case core.ModeRevisions:
		// From revisions, can only go to revisions (preserve resource group and app selection)
		return []list.Item{
			simpleContextItem{
				id:      "revisions",
				display: "🔄 App Revisions",
				enabled: true,
			},
		}

	case core.ModeContainers:
		// From containers, can only go to containers (preserve resource group, app, and revision selection)
		return []list.Item{
			simpleContextItem{
				id:      "containers",
				display: "🐳 Containers",
				enabled: true,
			},
		}

	case core.ModeEnvVars:
		// From env vars, can only go to env vars (preserve all selections)
		return []list.Item{
			simpleContextItem{
				id:      "env-vars",
				display: "🔧 Environment Variables",
				enabled: true,
			},
		}

	default:
		// Fallback - show top-level contexts
		return []list.Item{
			simpleContextItem{
				id:      "apps",
				display: "📦 Container Apps",
				enabled: true,
			},
			simpleContextItem{
				id:      "jobs",
				display: "⚡ Container App Jobs",
				enabled: false,
			},
		}
	}
}

func (m model) handleContextListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "?":
		// Handle help toggle in context mode
		m.help.ShowAll = !m.help.ShowAll
		// Recalculate context list size since help bar height changed
		m.core.SetContextList(m.createContextList())
		return m, nil

	case "enter":
		contextList := m.core.GetContextList()
		selectedItem, ok := contextList.SelectedItem().(simpleContextItem)
		if !ok {
			m.core.SetShowContextList(false)
			return m, nil
		}

		// Check if context is enabled
		if !selectedItem.enabled {
			m.core.SetStatusLine(selectedItem.display + " is not yet implemented")
			// Keep list open so user can select another option
			return m, nil
		}

		// Handle different types of selections
		var cmd tea.Cmd
		switch selectedItem.id {
		case "apps":
			// Navigate to container apps
			if m.core.GetCurrentMode() == core.ModeResourceGroups {
				// From resource groups to apps - need to select a resource group first
				m.core.SetStatusLine("Please select a resource group first")
				// Keep list open
				return m, nil
			} else {
				// From apps mode - stay in apps (preserve resource group selection)
				m.core.SetStatusLine("Container Apps")
			}

		case "jobs":
			// Context switching to jobs (not implemented yet)
			m.core.SetStatusLine(selectedItem.display + " is not yet implemented")
			// Keep list open
			return m, nil

		case "revisions":
			// Stay in revisions mode (preserve resource group and app selection)
			m.core.SetStatusLine("App Revisions")

		case "containers":
			// Stay in containers mode (preserve resource group, app, and revision selection)
			m.core.SetStatusLine("Containers")

		case "env-vars":
			// Stay in env vars mode (preserve all selections)
			m.core.SetStatusLine("Environment Variables")
		}

		m.core.SetShowContextList(false)
		return m, cmd

	case "esc", "q", ":":
		// Close without changing context
		m.core.SetShowContextList(false)
		return m, nil

	case "ctrl+c":
		// Allow quitting from context list
		return m, tea.Quit

	default:
		// Pass other keys to list component
		var cmd tea.Cmd
		contextList := m.core.GetContextList()
		contextList, cmd = contextList.Update(msg)
		m.core.SetContextList(contextList)
		return m, cmd
	}
}

func (m model) viewContextList() string {
	// Just show the list view - no header, no current state
	listView := m.core.GetContextList().View()

	// Create the context layout using the layout system
	layoutSystem := m.core.GetLayoutSystem()
	statusContext := layouts.StatusContext{
		Mode:          m.core.GetModeString(),
		StatusMessage: m.core.GetStatusLine(),
		Error:         m.core.GetError(),
		ContextInfo:   map[string]string{"breadcrumb": m.core.GetBreadcrumb()},
	}
	helpContext := layouts.HelpContext{
		Mode:    m.core.GetModeString(),
		ShowAll: m.help.ShowAll,
	}

	// If help is showing, pre-render the Bubble Tea help for context list
	if m.help.ShowAll {
		keyMap := m.getContextKeyMap()
		helpContext.BubbleTeaHelp = m.help.View(keyMap)
	}

	return layoutSystem.CreateContextLayout(listView, statusContext, helpContext)
}

// getContextKeyMap returns a key map based on the current page context
func (m model) getContextKeyMap() help.KeyMap {
	// Get the current page's help keys
	currentPage := m.core.GetCurrentPage()

	// Try to get help keys from the current page
	if helpKeysProvider, ok := currentPage.(interface{ GetHelpKeys() []key.Binding }); ok {
		helpKeys := helpKeysProvider.GetHelpKeys()
		return &contextKeyMap{keys: helpKeys}
	}

	// Fallback to main model keys
	return m.keys
}

// contextKeyMap wraps page-specific keys to implement help.KeyMap
type contextKeyMap struct {
	keys []key.Binding
}

func (ckm *contextKeyMap) ShortHelp() []key.Binding {
	if len(ckm.keys) <= 2 {
		return ckm.keys
	}
	return ckm.keys[:2]
}

func (ckm *contextKeyMap) FullHelp() [][]key.Binding {
	if len(ckm.keys) == 0 {
		return [][]key.Binding{}
	}

	// Group keys into rows of 4 for better display
	var rows [][]key.Binding
	for i := 0; i < len(ckm.keys); i += 4 {
		end := i + 4
		if end > len(ckm.keys) {
			end = len(ckm.keys)
		}
		rows = append(rows, ckm.keys[i:end])
	}
	return rows
}

// Helper methods

func (m model) withConfirm(text string, onYes func(model) (model, tea.Cmd), onNo func(model) (model, tea.Cmd)) model {
	m.confirm.Visible = true
	m.confirm.Text = text
	m.confirm.OnYes = onYes
	m.confirm.OnNo = onNo
	return m
}

// createCommandProvider creates the appropriate command provider based on mock mode
func createCommandProvider(useMockMode bool) providers.CommandProvider {
	if useMockMode {
		return providers.NewMockCommandProvider()
	}
	return providers.NewAzureCommandProvider()
}
