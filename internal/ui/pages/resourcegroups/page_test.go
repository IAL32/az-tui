package resourcegroups

import (
	"fmt"
	"testing"

	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	tea "github.com/charmbracelet/bubbletea"
)

// Simple test data
func createTestResourceGroups() []models.ResourceGroup {
	return []models.ResourceGroup{
		{
			Name:     "rg-production-eastus",
			Location: "East US",
			State:    "Succeeded",
			Tags: map[string]string{
				"environment": "production",
				"team":        "platform",
			},
		},
		{
			Name:     "rg-staging-westus",
			Location: "West US",
			State:    "Succeeded",
			Tags: map[string]string{
				"environment": "staging",
				"team":        "development",
			},
		},
	}
}

// Test page creation and basic functionality
func TestResourceGroupsPageBasics(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewResourceGroupsPage(layoutSystem)

	if page == nil {
		t.Fatal("Failed to create ResourceGroups page")
	}

	// Test initial state
	if page.IsLoading() {
		t.Error("Page should not be loading initially")
	}

	if len(page.GetData()) != 0 {
		t.Error("Page should have no data initially")
	}

	// Test data loading
	testData := createTestResourceGroups()
	page.SetData(testData)

	if len(page.GetData()) != len(testData) {
		t.Errorf("Expected %d resource groups, got %d", len(testData), len(page.GetData()))
	}

	// Test table creation
	table := page.GetTable()
	if table.TotalRows() != len(testData) {
		t.Errorf("Table should have %d rows, got %d", len(testData), table.TotalRows())
	}
}

// Test tag formatting business logic
func TestTagFormatting(t *testing.T) {
	tests := []struct {
		name     string
		tags     map[string]string
		expected string
	}{
		{
			name:     "empty tags",
			tags:     map[string]string{},
			expected: "-",
		},
		{
			name:     "nil tags",
			tags:     nil,
			expected: "-",
		},
		{
			name: "single tag",
			tags: map[string]string{
				"environment": "production",
			},
			expected: "environment=production",
		},
		{
			name: "multiple tags",
			tags: map[string]string{
				"environment": "production",
				"team":        "platform",
			},
			// Tags are sorted, so we expect this order
			expected: "environment=production, team=platform",
		},
		{
			name: "empty value",
			tags: map[string]string{
				"empty": "",
			},
			expected: `empty=""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatResourceGroupTags(tt.tags)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// Test key handling
func TestKeyHandling(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewResourceGroupsPage(layoutSystem)
	page.SetData(createTestResourceGroups())

	tests := []struct {
		name         string
		key          tea.KeyMsg
		shouldHandle bool
	}{
		{
			name: "filter key",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("/"),
			},
			shouldHandle: true,
		},
		{
			name: "quit key",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("q"),
			},
			shouldHandle: true,
		},
		{
			name: "enter key",
			key: tea.KeyMsg{
				Type: tea.KeyEnter,
			},
			shouldHandle: true,
		},
		{
			name: "help key (should not handle)",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("?"),
			},
			shouldHandle: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, handled := page.HandleKeyMsg(tt.key)
			if handled != tt.shouldHandle {
				t.Errorf("Expected handled=%v, got %v", tt.shouldHandle, handled)
			}
		})
	}
}

// Test view rendering in different states
func TestViewRendering(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewResourceGroupsPage(layoutSystem)

	// Test loading state
	page.SetLoading(true)
	view := page.View()
	if view == "" {
		t.Error("Loading view should not be empty")
	}

	// Test error state
	page.SetLoading(false)
	page.SetError(fmt.Errorf("test error"))
	view = page.View()
	if view == "" {
		t.Error("Error view should not be empty")
	}

	// Test normal state with data
	page.SetError(nil)
	page.SetData(createTestResourceGroups())
	view = page.View()
	if view == "" {
		t.Error("Normal view should not be empty")
	}
}

// Test navigation function setup
func TestNavigationSetup(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewResourceGroupsPage(layoutSystem)

	// Test setting navigation function
	page.SetNavigateToAppsFunc(func(rg models.ResourceGroup) tea.Cmd {
		return nil
	})

	// Verify function was set (we can't easily test execution without complex setup)
	if page.navigateToAppsFunc == nil {
		t.Error("Navigation function should be set")
	}
}
