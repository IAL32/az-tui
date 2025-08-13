package layouts

import (
	"github.com/charmbracelet/lipgloss"
)

// Manager handles all layout operations and provides consistent UI structure
type Manager struct {
	// Terminal dimensions
	termW int
	termH int

	// Layout configuration
	config Config

	// Component factories
	statusBarFactory StatusBarFactory
	helpBarFactory   HelpBarFactory
	themeManager     *ThemeManager
}

// NewManager creates a new layout manager with the given configuration
func NewManager(termW, termH int, config Config) *Manager {
	return &Manager{
		termW:        termW,
		termH:        termH,
		config:       config,
		themeManager: NewThemeManager(config.Theme),
	}
}

// SetDimensions updates the terminal dimensions
func (m *Manager) SetDimensions(termW, termH int) {
	m.termW = termW
	m.termH = termH
}

// SetStatusBarFactory sets the status bar factory
func (m *Manager) SetStatusBarFactory(factory StatusBarFactory) {
	m.statusBarFactory = factory
}

// SetHelpBarFactory sets the help bar factory
func (m *Manager) SetHelpBarFactory(factory HelpBarFactory) {
	m.helpBarFactory = factory
}

// CreateLayout creates the main layout with content and bars
func (m *Manager) CreateLayout(content string, options LayoutOptions) string {
	// Create the bottom bars
	helpBar := ""
	statusBar := ""

	if m.helpBarFactory != nil {
		helpBar = m.helpBarFactory.CreateHelpBar(options.HelpContext)
	}

	if m.statusBarFactory != nil {
		statusBar = m.statusBarFactory.CreateStatusBar(options.StatusContext)
	}

	// Calculate available height for main content
	helpBarHeight := lipgloss.Height(helpBar)
	statusBarHeight := lipgloss.Height(statusBar)
	mainContentHeight := m.termH - helpBarHeight - statusBarHeight

	// Apply content styling if specified
	if options.ContentStyle != nil {
		content = options.ContentStyle.Render(content)
	}

	// Position main content at top, help bar and status bar at bottom
	body := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Height(mainContentHeight).Render(content),
		helpBar,
		statusBar,
	)

	// Handle modal overlay if present
	if options.Modal != nil && options.Modal.Visible {
		return m.CreateModalLayout(body, *options.Modal)
	}

	return body
}

// CreateModalLayout creates a modal overlay on top of the base layout
func (m *Manager) CreateModalLayout(baseLayout string, modal ModalOptions) string {
	modalContent := modal.Content
	if modal.Style != nil {
		modalContent = modal.Style.Render(modalContent)
	}

	return lipgloss.Place(
		m.termW,
		m.termH,
		lipgloss.Center,
		lipgloss.Center,
		modalContent,
	)
}

// GetContentDimensions returns the available dimensions for content
func (m *Manager) GetContentDimensions(options LayoutOptions) (width, height int) {
	helpBarHeight := 0
	statusBarHeight := 0

	if m.helpBarFactory != nil {
		helpBar := m.helpBarFactory.CreateHelpBar(options.HelpContext)
		helpBarHeight = lipgloss.Height(helpBar)
	}

	if m.statusBarFactory != nil {
		statusBar := m.statusBarFactory.CreateStatusBar(options.StatusContext)
		statusBarHeight = lipgloss.Height(statusBar)
	}

	return m.termW, m.termH - helpBarHeight - statusBarHeight
}

// GetTheme returns the current theme manager
func (m *Manager) GetTheme() *ThemeManager {
	return m.themeManager
}

// UpdateTheme updates the theme configuration
func (m *Manager) UpdateTheme(theme ThemeConfig) {
	m.themeManager = NewThemeManager(theme)
	m.config.Theme = theme
}
