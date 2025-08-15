package layouts

import (
	"github.com/charmbracelet/lipgloss"
)

// TemplateManager manages layout templates
type TemplateManager struct {
	templates map[string]LayoutTemplate
	manager   *Manager
}

// NewTemplateManager creates a new template manager
func NewTemplateManager(manager *Manager) *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[string]LayoutTemplate),
		manager:   manager,
	}
	tm.registerDefaultTemplates()
	return tm
}

// registerDefaultTemplates registers the built-in layout templates
func (tm *TemplateManager) registerDefaultTemplates() {
	// Table layout template
	tm.RegisterTemplate(LayoutTemplate{
		Name:        "table",
		Description: "Standard table layout with status and help bars",
		Options: LayoutOptions{
			FillHeight: true,
		},
		ContentFunc: tm.createTableLayout,
	})

	// Loading layout template
	tm.RegisterTemplate(LayoutTemplate{
		Name:        "loading",
		Description: "Loading state layout with centered spinner and message",
		Options: LayoutOptions{
			CenterContent: true,
		},
		ContentFunc: tm.createLoadingLayout,
	})

	// Error layout template
	tm.RegisterTemplate(LayoutTemplate{
		Name:        "error",
		Description: "Error state layout with error message and help text",
		Options: LayoutOptions{
			CenterContent: true,
		},
		ContentFunc: tm.createErrorLayout,
	})

	// Modal layout template
	tm.RegisterTemplate(LayoutTemplate{
		Name:        "modal",
		Description: "Modal overlay layout for confirmations and dialogs",
		Options: LayoutOptions{
			Modal: &ModalOptions{
				Visible: true,
			},
		},
		ContentFunc: tm.createModalLayout,
	})

	// Context layout template
	tm.RegisterTemplate(LayoutTemplate{
		Name:        "context",
		Description: "Context selection layout for navigation",
		Options: LayoutOptions{
			FillHeight: true,
		},
		ContentFunc: tm.createContextLayout,
	})
}

// RegisterTemplate registers a new layout template
func (tm *TemplateManager) RegisterTemplate(template LayoutTemplate) {
	tm.templates[template.Name] = template
}

// GetTemplate returns a template by name
func (tm *TemplateManager) GetTemplate(name string) (*LayoutTemplate, bool) {
	template, exists := tm.templates[name]
	return &template, exists
}

// ListTemplates returns all available template names
func (tm *TemplateManager) ListTemplates() []string {
	names := make([]string, 0, len(tm.templates))
	for name := range tm.templates {
		names = append(names, name)
	}
	return names
}

// RenderTemplate renders content using the specified template
func (tm *TemplateManager) RenderTemplate(templateName, content string, options LayoutOptions) string {
	template, exists := tm.GetTemplate(templateName)
	if !exists {
		// Fallback to basic layout
		return tm.manager.CreateLayout(content, options)
	}

	// Merge template options with provided options
	mergedOptions := tm.mergeOptions(template.Options, options)

	// Use template's content function if available
	if template.ContentFunc != nil {
		state := LayoutState{
			TerminalWidth:  tm.manager.termW,
			TerminalHeight: tm.manager.termH,
		}
		content = template.ContentFunc(content, state)
	}

	return tm.manager.CreateLayout(content, mergedOptions)
}

// mergeOptions merges template options with provided options
func (tm *TemplateManager) mergeOptions(templateOpts, providedOpts LayoutOptions) LayoutOptions {
	merged := templateOpts

	// Override with provided options
	if providedOpts.ContentStyle != nil {
		merged.ContentStyle = providedOpts.ContentStyle
	}
	if providedOpts.Modal != nil {
		merged.Modal = providedOpts.Modal
	}
	// Always use provided help context if it has any meaningful content
	// Note: Mode 0 (ModeResourceGroups) is valid, so we can't use Mode != 0 as a check
	merged.HelpContext = providedOpts.HelpContext

	// Always use provided status context if it has any meaningful content
	// Note: Mode 0 (ModeResourceGroups) is valid, so we can't use Mode != 0 as a check
	merged.StatusContext = providedOpts.StatusContext

	return merged
}

// Template content functions

// createTableLayout creates a table layout
func (tm *TemplateManager) createTableLayout(content string, state LayoutState) string {
	return content // Table content is used as-is
}

// createLoadingLayout creates a loading layout with centered content
func (tm *TemplateManager) createLoadingLayout(content string, state LayoutState) string {
	theme := tm.manager.GetTheme()

	loadingContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		theme.GetStyle("accent").Render(content),
		"",
	)

	return loadingContent
}

// createErrorLayout creates an error layout with error message and help
func (tm *TemplateManager) createErrorLayout(content string, state LayoutState) string {
	theme := tm.manager.GetTheme()

	// Parse content as "error|help" format
	errorMsg := content
	helpMsg := "Press r to retry or q to quit."

	// Try to split content if it contains separator
	if parts := lipgloss.NewStyle().Render(content); len(parts) > 0 {
		// For now, use content as error message
		errorMsg = content
	}

	errorContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		theme.GetStyle("error").Render("Error: ")+errorMsg,
		theme.GetStyle("accent").Render(helpMsg),
		"",
	)

	return errorContent
}

// createModalLayout creates a modal layout
func (tm *TemplateManager) createModalLayout(content string, state LayoutState) string {
	theme := tm.manager.GetTheme()

	// Create modal box styling
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.GetStyle("accent").GetForeground()).
		Padding(1, 2).
		Background(theme.GetTheme().Colors.Background).
		Foreground(theme.GetTheme().Colors.Foreground)

	return modalStyle.Render(content)
}

// createContextLayout creates a context selection layout
func (tm *TemplateManager) createContextLayout(content string, state LayoutState) string {
	return content // Context content is used as-is
}

// Convenience methods for common layouts

// CreateTableLayout creates a table layout
func (tm *TemplateManager) CreateTableLayout(tableView string, options LayoutOptions) string {
	return tm.RenderTemplate("table", tableView, options)
}

// CreateLoadingLayout creates a loading layout
func (tm *TemplateManager) CreateLoadingLayout(message string, options LayoutOptions) string {
	return tm.RenderTemplate("loading", message, options)
}

// CreateErrorLayout creates an error layout
func (tm *TemplateManager) CreateErrorLayout(errorMsg, helpMsg string, options LayoutOptions) string {
	content := errorMsg
	if helpMsg != "" {
		content = errorMsg + "|" + helpMsg
	}
	return tm.RenderTemplate("error", content, options)
}

// CreateModalLayout creates a modal layout
func (tm *TemplateManager) CreateModalLayout(content string, options LayoutOptions) string {
	if options.Modal == nil {
		options.Modal = &ModalOptions{Visible: true}
	}
	options.Modal.Content = content
	return tm.RenderTemplate("modal", content, options)
}

// CreateContextLayout creates a context selection layout
func (tm *TemplateManager) CreateContextLayout(listView string, options LayoutOptions) string {
	return tm.RenderTemplate("context", listView, options)
}
