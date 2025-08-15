package layouts

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// DefaultStatusBarFactory provides the default status bar implementation
type DefaultStatusBarFactory struct {
	theme *ThemeManager
	termW int
}

// NewDefaultStatusBarFactory creates a new default status bar factory
func NewDefaultStatusBarFactory(theme *ThemeManager, termW int) *DefaultStatusBarFactory {
	return &DefaultStatusBarFactory{
		theme: theme,
		termW: termW,
	}
}

// CreateStatusBar creates a status bar based on the provided context
func (f *DefaultStatusBarFactory) CreateStatusBar(context StatusContext) string {
	w := lipgloss.Width

	// Mode indicator
	var modeIndicator string
	switch context.Mode {
	case ModeApps:
		modeIndicator = f.theme.GetStyle("modeApps").Render("📦 APPS")
	case ModeRevisions:
		modeIndicator = f.theme.GetStyle("modeRevisions").Render("🔄 REVISIONS")
	case ModeContainers:
		modeIndicator = f.theme.GetStyle("modeContainers").Render("🐳 CONTAINERS")
	case ModeEnvVars:
		modeIndicator = f.theme.GetStyle("modeContainers").Render("🔧 ENV VARS")
	case ModeResourceGroups:
		modeIndicator = f.theme.GetStyle("modeApps").Render("📁 RESOURCE GROUPS")
	default:
		modeIndicator = f.theme.GetStyle("modeApps").Render("📦 APPS")
	}

	// Status indicator
	var statusIndicator string
	if context.FilterActive {
		statusIndicator = f.theme.GetStyle("statusLoading").Render("Filtering")
	} else if context.Loading {
		statusIndicator = f.theme.GetStyle("statusLoading").Render("Loading")
	} else if context.Error != nil {
		statusIndicator = f.theme.GetStyle("statusError").Render("Error")
	} else {
		statusIndicator = f.theme.GetStyle("statusReady").Render("Ready")
	}

	// Count indicators
	var countIndicators []string
	for name, count := range context.Counters {
		var label string
		if count == 1 {
			label = fmt.Sprintf("1 %s", name)
		} else {
			label = fmt.Sprintf("%d %ss", count, name)
		}
		countIndicators = append(countIndicators, f.theme.GetStyle("count").Render(label))
	}

	// Context info indicators
	var contextIndicators []string
	// Define consistent key order to ensure deterministic display
	keyOrder := []string{"app", "revision", "container", "resource_group"}
	for _, name := range keyOrder {
		if value, exists := context.ContextInfo[name]; exists {
			indicator := f.theme.GetStyle("context").Render(fmt.Sprintf("%s: %s", name, value))
			contextIndicators = append(contextIndicators, indicator)
		}
	}

	// Calculate widths for fixed elements
	fixedWidth := w(modeIndicator) + w(statusIndicator)
	for _, indicator := range countIndicators {
		fixedWidth += w(indicator)
	}
	for _, indicator := range contextIndicators {
		fixedWidth += w(indicator)
	}

	// Status message (expandable middle section)
	statusMessage := context.StatusMessage
	if statusMessage == "" {
		if context.Error != nil {
			statusMessage = context.Error.Error()
		} else if context.Loading {
			statusMessage = "Loading..."
		} else {
			statusMessage = ""
		}
	}

	// Create expandable status text
	statusVal := f.theme.GetStyle("statusText").
		Width(max(0, f.termW-fixedWidth-4)). // Leave some margin
		Render(statusMessage)

	// Build the status bar
	var elements []string
	elements = append(elements, modeIndicator)
	elements = append(elements, statusIndicator)
	elements = append(elements, contextIndicators...)
	elements = append(elements, statusVal)
	elements = append(elements, countIndicators...)

	bar := lipgloss.JoinHorizontal(lipgloss.Top, elements...)
	return f.theme.GetStyle("statusBar").Width(f.termW).Render(bar)
}

// SetTerminalWidth updates the terminal width
func (f *DefaultStatusBarFactory) SetTerminalWidth(width int) {
	f.termW = width
}

// DefaultHelpBarFactory provides the default help bar implementation
type DefaultHelpBarFactory struct {
	theme *ThemeManager
}

// NewDefaultHelpBarFactory creates a new default help bar factory
func NewDefaultHelpBarFactory(theme *ThemeManager) *DefaultHelpBarFactory {
	return &DefaultHelpBarFactory{
		theme: theme,
	}
}

// CreateHelpBar creates a help bar based on the provided context
func (f *DefaultHelpBarFactory) CreateHelpBar(context HelpContext) string {
	// Build help text based on context
	var helpItems []string

	// Add mode-specific help
	switch context.Mode {
	case ModeApps:
		helpItems = append(helpItems, "enter: view revisions", "l: logs", "s/e: exec", "r: refresh", "/: filter", "esc: back", "?: help", "q: quit")
	case ModeRevisions:
		helpItems = append(helpItems, "enter: view containers", "R: restart", "l: logs", "s: exec", "r: refresh", "/: filter", "esc: back", "?: help", "q: quit")
	case ModeContainers:
		helpItems = append(helpItems, "v: env vars", "s: shell", "l: logs", "r: refresh", "/: filter", "esc: back", "?: help", "q: quit")
	case ModeEnvVars:
		helpItems = append(helpItems, "/: filter", "shift+←/→: scroll", "esc: back", "?: help", "q: quit")
	case ModeResourceGroups:
		helpItems = append(helpItems, "enter: select", "r: refresh", "/: filter", "?: help", "q: quit")
	default:
		helpItems = append(helpItems, "?: help", "q: quit")
	}

	// Add custom help items
	for key, desc := range context.CustomHelp {
		helpItems = append(helpItems, fmt.Sprintf("%s: %s", key, desc))
	}

	// Create help text
	helpText := ""
	if context.ShowAll {
		// When ShowAll is true, prefer Bubble Tea help content if provided
		if context.BubbleTeaHelp != "" {
			return context.BubbleTeaHelp
		}
		// Fallback to showing all help items when no Bubble Tea help provided
		if len(helpItems) > 0 {
			helpText = ""
			for i, item := range helpItems {
				if i > 0 {
					helpText += " • "
				}
				helpText += item
			}
		}
	} else {
		// Show minimal help
		if len(helpItems) > 0 {
			helpText = helpItems[0]
			if len(helpItems) > 1 {
				helpText += " • " + helpItems[len(helpItems)-1] // Show first and last
			}
		}
	}

	return helpText
}

// DefaultComponentFactory provides default component implementations
type DefaultComponentFactory struct {
	theme *ThemeManager
}

// NewDefaultComponentFactory creates a new default component factory
func NewDefaultComponentFactory(theme *ThemeManager) *DefaultComponentFactory {
	return &DefaultComponentFactory{
		theme: theme,
	}
}

// CreateLoadingComponent creates a loading component
func (f *DefaultComponentFactory) CreateLoadingComponent(message string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		f.theme.GetStyle("accent").Render(message),
		"",
	)
}

// CreateErrorComponent creates an error component
func (f *DefaultComponentFactory) CreateErrorComponent(error string, helpText string) string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		f.theme.GetStyle("error").Render("Error: ")+error,
		f.theme.GetStyle("accent").Render(helpText),
		"",
	)
}

// CreateTableComponent creates a table component wrapper
func (f *DefaultComponentFactory) CreateTableComponent(content string) string {
	return content // Table content is used as-is
}

// CreateModalComponent creates a modal component
func (f *DefaultComponentFactory) CreateModalComponent(content string, options ModalOptions) string {
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(f.theme.GetTheme().Colors.Border).
		Padding(1, 2).
		Background(f.theme.GetTheme().Colors.Background).
		Foreground(f.theme.GetTheme().Colors.Foreground)

	if options.Style != nil {
		return options.Style.Render(content)
	}

	return modalStyle.Render(content)
}
