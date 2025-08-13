// Package table provides reusable table components and utilities for the UI.
package table

import (
	"github.com/charmbracelet/bubbles/textinput"
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

// UnifiedTableConfig holds configuration for creating unified tables
type UnifiedTableConfig struct {
	Columns     []table.Column
	Rows        []table.Row
	FilterInput interface{}
	BaseStyle   lipgloss.Style
	MaxWidth    int
	MaxHeight   int
	FilterFunc  func([]table.Column) func(table.Row, string) bool
}

// CreateUnifiedTable creates a table with unified styling and common configuration
func CreateUnifiedTable(config UnifiedTableConfig) table.Model {
	t := table.New(config.Columns).
		WithRows(config.Rows).
		BorderRounded().
		WithBaseStyle(config.BaseStyle).
		WithMaxTotalWidth(config.MaxWidth).
		WithHorizontalFreezeColumnCount(1).
		Filtered(true).
		Focused(true)

	// Set filter function if provided
	if config.FilterFunc != nil {
		t = t.WithFilterFunc(config.FilterFunc(config.Columns))
	}

	// Set the filter input with proper type assertion
	if config.FilterInput != nil {
		if filter, ok := config.FilterInput.(textinput.Model); ok {
			t = t.WithFilterInput(filter)
		}
	}

	// Set page size if max height is specified
	if config.MaxHeight > 0 {
		t = t.WithPageSize(config.MaxHeight)
	}

	return t
}
