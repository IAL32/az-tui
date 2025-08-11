package ui

import (
	models "az-tui/internal/models"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// ColumnConfig defines configuration for a dynamic column
type ColumnConfig struct {
	Key        string
	Title      string
	MinWidth   int
	Padding    int
	WithFilter bool
}

// DynamicColumnBuilder helps create columns with dynamic widths based on content
type DynamicColumnBuilder struct {
	configs []ColumnConfig
	widths  map[string]int
}

// NewDynamicColumnBuilder creates a new dynamic column builder
func NewDynamicColumnBuilder() *DynamicColumnBuilder {
	return &DynamicColumnBuilder{
		configs: make([]ColumnConfig, 0),
		widths:  make(map[string]int),
	}
}

// AddColumn adds a column configuration to the builder
func (dcb *DynamicColumnBuilder) AddColumn(key, title string, minWidth int, withFilter bool) *DynamicColumnBuilder {
	dcb.configs = append(dcb.configs, ColumnConfig{
		Key:        key,
		Title:      title,
		MinWidth:   minWidth,
		Padding:    2, // Default padding
		WithFilter: withFilter,
	})
	// Initialize width with title length
	dcb.widths[key] = len(title)
	return dcb
}

// AddColumnWithPadding adds a column configuration with custom padding
func (dcb *DynamicColumnBuilder) AddColumnWithPadding(key, title string, minWidth, padding int, withFilter bool) *DynamicColumnBuilder {
	dcb.configs = append(dcb.configs, ColumnConfig{
		Key:        key,
		Title:      title,
		MinWidth:   minWidth,
		Padding:    padding,
		WithFilter: withFilter,
	})
	// Initialize width with title length
	dcb.widths[key] = len(title)
	return dcb
}

// UpdateWidth updates the width for a column based on content length
func (dcb *DynamicColumnBuilder) UpdateWidth(key string, contentLength int) {
	if currentWidth, exists := dcb.widths[key]; exists {
		if contentLength > currentWidth {
			dcb.widths[key] = contentLength
		}
	}
}

// UpdateWidthFromString updates the width for a column based on string content
func (dcb *DynamicColumnBuilder) UpdateWidthFromString(key, content string) {
	dcb.UpdateWidth(key, len(content))
}

// Build creates the final table columns with calculated widths
func (dcb *DynamicColumnBuilder) Build() []table.Column {
	columns := make([]table.Column, len(dcb.configs))

	for i, config := range dcb.configs {
		// Calculate final width: content width + padding, with minimum width enforced
		finalWidth := dcb.widths[config.Key] + config.Padding
		if finalWidth < config.MinWidth {
			finalWidth = config.MinWidth
		}

		// Create column
		column := table.NewColumn(config.Key, config.Title, finalWidth)
		if config.WithFilter {
			column = column.WithFiltered(true)
		}

		columns[i] = column
	}

	return columns
}

// Helper methods for component factories

func (m model) getCurrentAppName() string {
	if app, ok := m.currentApp(); ok {
		return app.Name
	}
	return ""
}

func (m model) getCurrentApp() models.ContainerApp {
	if app, ok := m.currentApp(); ok {
		return app
	}
	return models.ContainerApp{}
}

// Confirmation dialog helper
func (m model) confirmBox() string {
	if !m.confirm.Visible {
		return ""
	}

	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Padding(1, 2).
		Width(50).
		Align(lipgloss.Center).
		Render(m.confirm.Text + "\n\n[y] Yes  [n] No")

	return box
}

// getStatusColor returns the appropriate color for a status value
func getStatusColor(status string) string {
	switch strings.ToLower(status) {
	case "running", "succeeded", "healthy", "ready", "active":
		return "#8c8" // Green
	case "failed", "error", "unhealthy", "critical":
		return "#c88" // Red
	case "pending", "provisioning", "starting", "updating":
		return "#cc8" // Yellow
	case "stopped", "inactive", "unknown", "-":
		return "#888" // Gray
	default:
		return "#888" // Default gray
	}
}
