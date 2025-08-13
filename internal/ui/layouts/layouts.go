package layouts

import (
	"github.com/charmbracelet/lipgloss"
)

// LayoutSystem provides a complete layout management system
type LayoutSystem struct {
	manager          *Manager
	templateManager  *TemplateManager
	statusFactory    StatusBarFactory
	helpFactory      HelpBarFactory
	componentFactory ComponentFactory
}

// NewLayoutSystem creates a new layout system with default configuration
func NewLayoutSystem(termW, termH int) *LayoutSystem {
	config := DefaultConfig()
	manager := NewManager(termW, termH, config)

	// Create default factories
	theme := manager.GetTheme()
	statusFactory := NewDefaultStatusBarFactory(theme, termW)
	helpFactory := NewDefaultHelpBarFactory(theme)
	componentFactory := NewDefaultComponentFactory(theme)

	// Set factories in manager
	manager.SetStatusBarFactory(statusFactory)
	manager.SetHelpBarFactory(helpFactory)

	// Create template manager
	templateManager := NewTemplateManager(manager)

	return &LayoutSystem{
		manager:          manager,
		templateManager:  templateManager,
		statusFactory:    statusFactory,
		helpFactory:      helpFactory,
		componentFactory: componentFactory,
	}
}

// NewLayoutSystemWithConfig creates a new layout system with custom configuration
func NewLayoutSystemWithConfig(termW, termH int, config Config) *LayoutSystem {
	manager := NewManager(termW, termH, config)

	// Create default factories
	theme := manager.GetTheme()
	statusFactory := NewDefaultStatusBarFactory(theme, termW)
	helpFactory := NewDefaultHelpBarFactory(theme)
	componentFactory := NewDefaultComponentFactory(theme)

	// Set factories in manager
	manager.SetStatusBarFactory(statusFactory)
	manager.SetHelpBarFactory(helpFactory)

	// Create template manager
	templateManager := NewTemplateManager(manager)

	return &LayoutSystem{
		manager:          manager,
		templateManager:  templateManager,
		statusFactory:    statusFactory,
		helpFactory:      helpFactory,
		componentFactory: componentFactory,
	}
}

// SetDimensions updates the terminal dimensions
func (ls *LayoutSystem) SetDimensions(termW, termH int) {
	ls.manager.SetDimensions(termW, termH)
	if factory, ok := ls.statusFactory.(*DefaultStatusBarFactory); ok {
		factory.SetTerminalWidth(termW)
	}
}

// Layout creation methods

// CreateLayout creates a basic layout with content and bars
func (ls *LayoutSystem) CreateLayout(content string, options LayoutOptions) string {
	return ls.manager.CreateLayout(content, options)
}

// CreateTableLayout creates a table layout
func (ls *LayoutSystem) CreateTableLayout(tableView string, statusContext StatusContext, helpContext HelpContext) string {
	options := LayoutOptions{
		StatusContext: statusContext,
		HelpContext:   helpContext,
		FillHeight:    true,
	}
	return ls.templateManager.CreateTableLayout(tableView, options)
}

// CreateLoadingLayout creates a loading layout
func (ls *LayoutSystem) CreateLoadingLayout(message string, statusContext StatusContext, helpContext HelpContext) string {
	options := LayoutOptions{
		StatusContext: statusContext,
		HelpContext:   helpContext,
		CenterContent: true,
	}
	return ls.templateManager.CreateLoadingLayout(message, options)
}

// CreateErrorLayout creates an error layout
func (ls *LayoutSystem) CreateErrorLayout(errorMsg, helpMsg string, statusContext StatusContext, helpContext HelpContext) string {
	options := LayoutOptions{
		StatusContext: statusContext,
		HelpContext:   helpContext,
		CenterContent: true,
	}
	return ls.templateManager.CreateErrorLayout(errorMsg, helpMsg, options)
}

// CreateModalLayout creates a modal layout
func (ls *LayoutSystem) CreateModalLayout(content string, statusContext StatusContext, helpContext HelpContext) string {
	options := LayoutOptions{
		StatusContext: statusContext,
		HelpContext:   helpContext,
		Modal: &ModalOptions{
			Visible: true,
			Content: content,
		},
	}
	return ls.templateManager.CreateModalLayout(content, options)
}

// CreateContextLayout creates a context selection layout
func (ls *LayoutSystem) CreateContextLayout(listView string, statusContext StatusContext, helpContext HelpContext) string {
	options := LayoutOptions{
		StatusContext: statusContext,
		HelpContext:   helpContext,
		FillHeight:    true,
	}
	return ls.templateManager.CreateContextLayout(listView, options)
}

// Component creation methods

// CreateLoadingComponent creates a loading component
func (ls *LayoutSystem) CreateLoadingComponent(message string) string {
	return ls.componentFactory.CreateLoadingComponent(message)
}

// CreateErrorComponent creates an error component
func (ls *LayoutSystem) CreateErrorComponent(error string, helpText string) string {
	return ls.componentFactory.CreateErrorComponent(error, helpText)
}

// CreateTableComponent creates a table component wrapper
func (ls *LayoutSystem) CreateTableComponent(content string) string {
	return ls.componentFactory.CreateTableComponent(content)
}

// CreateModalComponent creates a modal component
func (ls *LayoutSystem) CreateModalComponent(content string, options ModalOptions) string {
	return ls.componentFactory.CreateModalComponent(content, options)
}

// Template management methods

// RegisterTemplate registers a custom layout template
func (ls *LayoutSystem) RegisterTemplate(template LayoutTemplate) {
	ls.templateManager.RegisterTemplate(template)
}

// GetTemplate returns a template by name
func (ls *LayoutSystem) GetTemplate(name string) (*LayoutTemplate, bool) {
	return ls.templateManager.GetTemplate(name)
}

// ListTemplates returns all available template names
func (ls *LayoutSystem) ListTemplates() []string {
	return ls.templateManager.ListTemplates()
}

// RenderTemplate renders content using a specific template
func (ls *LayoutSystem) RenderTemplate(templateName, content string, options LayoutOptions) string {
	return ls.templateManager.RenderTemplate(templateName, content, options)
}

// Theme management methods

// GetTheme returns the current theme manager
func (ls *LayoutSystem) GetTheme() *ThemeManager {
	return ls.manager.GetTheme()
}

// UpdateTheme updates the theme configuration
func (ls *LayoutSystem) UpdateTheme(theme ThemeConfig) {
	ls.manager.UpdateTheme(theme)

	// Update factories with new theme
	newTheme := ls.manager.GetTheme()
	if factory, ok := ls.statusFactory.(*DefaultStatusBarFactory); ok {
		factory.theme = newTheme
	}
	if factory, ok := ls.helpFactory.(*DefaultHelpBarFactory); ok {
		factory.theme = newTheme
	}
	if factory, ok := ls.componentFactory.(*DefaultComponentFactory); ok {
		factory.theme = newTheme
	}
}

// GetStyle returns a style by name from the theme
func (ls *LayoutSystem) GetStyle(name string) lipgloss.Style {
	return ls.manager.GetTheme().GetStyle(name)
}

// Utility methods

// GetContentDimensions returns the available dimensions for content
func (ls *LayoutSystem) GetContentDimensions(options LayoutOptions) (width, height int) {
	return ls.manager.GetContentDimensions(options)
}

// SetCustomStatusBarFactory sets a custom status bar factory
func (ls *LayoutSystem) SetCustomStatusBarFactory(factory StatusBarFactory) {
	ls.statusFactory = factory
	ls.manager.SetStatusBarFactory(factory)
}

// SetCustomHelpBarFactory sets a custom help bar factory
func (ls *LayoutSystem) SetCustomHelpBarFactory(factory HelpBarFactory) {
	ls.helpFactory = factory
	ls.manager.SetHelpBarFactory(factory)
}

// SetCustomComponentFactory sets a custom component factory
func (ls *LayoutSystem) SetCustomComponentFactory(factory ComponentFactory) {
	ls.componentFactory = factory
}

// Convenience functions for backward compatibility

// CreateBasicLayout creates a basic layout (backward compatibility)
func CreateBasicLayout(content string, termW, termH int) string {
	ls := NewLayoutSystem(termW, termH)
	return ls.CreateLayout(content, LayoutOptions{})
}

// CreateStyledLayout creates a layout with custom styling
func CreateStyledLayout(content string, style lipgloss.Style, termW, termH int) string {
	ls := NewLayoutSystem(termW, termH)
	options := LayoutOptions{
		ContentStyle: &style,
	}
	return ls.CreateLayout(content, options)
}
