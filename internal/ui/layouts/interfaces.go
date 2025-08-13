package layouts

import "github.com/charmbracelet/lipgloss"

// StatusBarFactory defines the interface for creating status bars
type StatusBarFactory interface {
	CreateStatusBar(context StatusContext) string
}

// HelpBarFactory defines the interface for creating help bars
type HelpBarFactory interface {
	CreateHelpBar(context HelpContext) string
}

// LayoutRenderer defines the interface for rendering layouts
type LayoutRenderer interface {
	RenderLayout(content string, options LayoutOptions) string
	RenderModal(content string, modal ModalOptions) string
	GetContentDimensions(options LayoutOptions) (width, height int)
}

// ThemeProvider defines the interface for theme management
type ThemeProvider interface {
	GetTheme() ThemeConfig
	SetTheme(theme ThemeConfig)
	GetStyle(name string) *lipgloss.Style
	UpdateStyle(name string, style lipgloss.Style)
}

// ComponentFactory defines the interface for creating UI components
type ComponentFactory interface {
	CreateLoadingComponent(message string) string
	CreateErrorComponent(error string, helpText string) string
	CreateTableComponent(content string) string
	CreateModalComponent(content string, options ModalOptions) string
}

// LayoutTemplate defines the interface for layout templates
type LayoutTemplateProvider interface {
	GetTemplate(name string) (*LayoutTemplate, bool)
	RegisterTemplate(template LayoutTemplate)
	ListTemplates() []string
}

// ResponsiveManager defines the interface for responsive layout management
type ResponsiveManager interface {
	GetBreakpoint(width int) ResponsiveBreakpoint
	UpdateLayout(width, height int) Config
	IsResponsive() bool
}
