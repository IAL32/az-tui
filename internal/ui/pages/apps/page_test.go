package apps

import (
	"fmt"
	"strings"
	"testing"

	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/ui/layouts"
	tea "github.com/charmbracelet/bubbletea"
)

// Simple test data
func createTestContainerApps() []models.ContainerApp {
	return []models.ContainerApp{
		{
			Name:              "web-app-prod",
			ResourceGroup:     "rg-production-eastus",
			Location:          "East US",
			LatestRevision:    "web-app-prod--rev1",
			IngressFQDN:       "web-app-prod.example.com",
			ProvisioningState: "Succeeded",
			RunningStatus:     "Running",
			MinReplicas:       2,
			MaxReplicas:       10,
			CPU:               0.5,
			Memory:            "1Gi",
			IngressExternal:   true,
			TargetPort:        80,
		},
		{
			Name:              "api-service-prod",
			ResourceGroup:     "rg-production-eastus",
			Location:          "East US",
			LatestRevision:    "api-service-prod--rev2",
			ProvisioningState: "Succeeded",
			RunningStatus:     "Running",
			MinReplicas:       1,
			MaxReplicas:       5,
			CPU:               1.0,
			Memory:            "2Gi",
		},
	}
}

// Test page creation and basic functionality
func TestAppsPageBasics(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewAppsPage(layoutSystem)

	if page == nil {
		t.Fatal("Failed to create Apps page")
	}

	// Test initial state
	if page.IsLoading() {
		t.Error("Page should not be loading initially")
	}

	if len(page.GetData()) != 0 {
		t.Error("Page should have no data initially")
	}

	// Test resource group context
	page.SetResourceGroupContext("rg-production-eastus")
	if page.resourceGroupName != "rg-production-eastus" {
		t.Error("Resource group context should be set")
	}

	// Test data loading
	testData := createTestContainerApps()
	page.SetData(testData)

	if len(page.GetData()) != len(testData) {
		t.Errorf("Expected %d container apps, got %d", len(testData), len(page.GetData()))
	}

	// Test table creation
	table := page.GetTable()
	if table.TotalRows() != len(testData) {
		t.Errorf("Table should have %d rows, got %d", len(testData), table.TotalRows())
	}
}

// Test action handling
func TestAppsPageActions(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewAppsPage(layoutSystem)
	page.SetData(createTestContainerApps())

	// Test logs action setup
	logsCalled := false
	page.SetShowLogsFunc(func(app models.ContainerApp) tea.Cmd {
		logsCalled = true
		return nil
	})

	// Test logs key
	logsMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("l"),
	}

	_, handled := page.HandleKeyMsg(logsMsg)
	if !handled {
		t.Error("Logs key 'l' should be handled")
	}

	// Test exec action setup
	execCalled := false
	page.SetExecIntoAppFunc(func(app models.ContainerApp) tea.Cmd {
		execCalled = true
		return nil
	})

	// Test exec key
	execMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("e"),
	}

	_, handled = page.HandleKeyMsg(execMsg)
	if !handled {
		t.Error("Exec key 'e' should be handled")
	}

	// Use variables to avoid "declared and not used" errors
	_ = logsCalled
	_ = execCalled
}

// Test navigation handling
func TestAppsPageNavigation(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewAppsPage(layoutSystem)
	page.SetData(createTestContainerApps())

	// Test forward navigation setup
	page.SetNavigateToRevisionsFunc(func(app models.ContainerApp) tea.Cmd {
		return nil
	})

	// Test Enter key
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	_, handled := page.HandleKeyMsg(enterMsg)
	if !handled {
		t.Error("Enter key should be handled")
	}

	// Test back navigation setup
	backCalled := false
	page.SetBackToResourceGroupsFunc(func() tea.Cmd {
		backCalled = true
		return nil
	})

	// Test ESC key
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	cmd, handled := page.HandleKeyMsg(escMsg)
	if !handled {
		t.Error("ESC key should be handled")
	}

	// Execute the command to trigger the callback
	if cmd != nil {
		cmd()
	}

	if !backCalled {
		t.Error("Back function should have been called")
	}
}

// Test key handling
func TestAppsPageKeyHandling(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewAppsPage(layoutSystem)
	page.SetData(createTestContainerApps())

	tests := []struct {
		name         string
		key          tea.KeyMsg
		shouldHandle bool
	}{
		{
			name: "logs key",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("l"),
			},
			shouldHandle: true,
		},
		{
			name: "exec key",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("e"),
			},
			shouldHandle: true,
		},
		{
			name: "alternative exec key",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("s"),
			},
			shouldHandle: true,
		},
		{
			name: "filter key",
			key: tea.KeyMsg{
				Type:  tea.KeyRunes,
				Runes: []rune("/"),
			},
			shouldHandle: true,
		},
		{
			name: "esc key",
			key: tea.KeyMsg{
				Type: tea.KeyEsc,
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
func TestAppsPageViewRendering(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewAppsPage(layoutSystem)
	page.SetResourceGroupContext("rg-production-eastus")

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
	page.SetData(createTestContainerApps())
	view = page.View()
	if view == "" {
		t.Error("Normal view should not be empty")
	}
}

// Test help keys configuration
func TestAppsPageHelpKeys(t *testing.T) {
	layoutSystem := layouts.NewLayoutSystem(80, 24)
	page := NewAppsPage(layoutSystem)

	helpKeys := page.GetHelpKeys()
	if len(helpKeys) == 0 {
		t.Error("Page should have help keys defined")
	}

	// Check that essential keys are present
	keyStrings := make([]string, len(helpKeys))
	for i, key := range helpKeys {
		keyStrings[i] = key.Help().Key + ":" + key.Help().Desc
	}

	hasLogsKey := false
	hasExecKey := false
	hasEnterKey := false
	hasBackKey := false

	for _, keyStr := range keyStrings {
		if strings.Contains(keyStr, "l") && strings.Contains(keyStr, "logs") {
			hasLogsKey = true
		}
		if (strings.Contains(keyStr, "s") || strings.Contains(keyStr, "e")) && strings.Contains(keyStr, "exec") {
			hasExecKey = true
		}
		if strings.Contains(keyStr, "enter") {
			hasEnterKey = true
		}
		if strings.Contains(keyStr, "esc") && strings.Contains(keyStr, "back") {
			hasBackKey = true
		}
	}

	if !hasLogsKey {
		t.Error("Help keys should include logs key")
	}
	if !hasExecKey {
		t.Error("Help keys should include exec key")
	}
	if !hasEnterKey {
		t.Error("Help keys should include enter key")
	}
	if !hasBackKey {
		t.Error("Help keys should include back key")
	}
}
