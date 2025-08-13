package ui

import (
	"testing"
	"time"

	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/ui/core"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// Helper functions for testing
func createTestModelWithData(t *testing.T) model {
	t.Helper()
	m := InitialModel(true)

	// Simulate that data has been loaded by setting loading to false through the core interface
	// In a real scenario, this would happen through the Init() and Update() cycle
	if m.core != nil {
		// The core system manages loading state internally
		// We can't directly set loading state, but we can simulate loaded data
		msg := core.LoadedResourceGroupsMsg{
			ResourceGroups: []models.ResourceGroup{},
			Error:          nil,
		}
		m.core.HandleMessage(msg)
	}

	return m
}

// TestBasicModelCreation tests that we can create and initialize a model
func TestBasicModelCreation(t *testing.T) {
	// Create model with mock mode enabled
	m := InitialModel(true)

	// Verify initial state using core interface
	if m.core != nil {
		if m.core.GetCurrentMode() != core.ModeResourceGroups {
			t.Errorf("Expected initial mode to be ModeResourceGroups, got %v", m.core.GetCurrentMode())
		}

		// Check if there's an error (which would indicate loading completed with error)
		// or if we can determine loading state through other means
		if err := m.core.GetError(); err != nil {
			t.Logf("Model has error state: %v", err)
		}
	}

	if m.termW != 80 || m.termH != 24 {
		t.Errorf("Expected default terminal size 80x24, got %dx%d", m.termW, m.termH)
	}
}

// TestModelUpdate tests that the model can handle basic updates
func TestModelUpdate(t *testing.T) {
	// Create model with mock mode enabled
	m := InitialModel(true)

	// Test window size message
	updatedTeaModel, _ := m.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	updatedModel := updatedTeaModel.(model)

	if updatedModel.termW != 100 || updatedModel.termH != 30 {
		t.Errorf("Expected terminal size 100x30, got %dx%d", updatedModel.termW, updatedModel.termH)
	}
}

// TestKeyHandling tests basic key handling
func TestKeyHandling(t *testing.T) {
	// Create model with mock mode enabled
	m := InitialModel(true)

	// Test quit key handling directly on the model
	quitMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	}

	// The model should handle the quit key without panicking
	_, cmd := m.Update(quitMsg)

	if cmd == nil {
		t.Errorf("Expected no command for quit key, got %v", cmd)
	}
}

// TestModeConstants tests that mode constants are properly defined
func TestModeConstants(t *testing.T) {
	modes := []core.Mode{
		core.ModeApps,
		core.ModeRevisions,
		core.ModeContainers,
		core.ModeEnvVars,
		core.ModeResourceGroups,
	}

	// Verify all modes are different
	modeSet := make(map[core.Mode]bool)
	for _, m := range modes {
		if modeSet[m] {
			t.Errorf("Duplicate mode value: %v", m)
		}
		modeSet[m] = true
	}

	// Verify we have the expected number of modes
	if len(modes) != 5 {
		t.Errorf("Expected 5 modes, got %d", len(modes))
	}
}

// TestKeyBindings tests that all key bindings are properly configured
func TestKeyBindings(t *testing.T) {
	m := InitialModel(true)

	// Test that key bindings are initialized
	if m.keys.Enter.Keys() == nil {
		t.Error("Enter key binding not initialized")
	}

	if m.keys.Back.Keys() == nil {
		t.Error("Back key binding not initialized")
	}

	if m.keys.Quit.Keys() == nil {
		t.Error("Quit key binding not initialized")
	}

	if m.keys.Refresh.Keys() == nil {
		t.Error("Refresh key binding not initialized")
	}

	// Test help system
	if len(m.keys.ShortHelp()) == 0 {
		t.Error("Short help should return key bindings")
	}

	if len(m.keys.FullHelp()) == 0 {
		t.Error("Full help should return key bindings")
	}
}

// TestSpinnerInitialization tests that the spinner is properly initialized
func TestSpinnerInitialization(t *testing.T) {
	m := InitialModel(true)

	// Test that spinner is initialized
	spinnerView := m.spinner.View()
	if spinnerView == "" {
		t.Error("Spinner should have a view")
	}

	// Test spinner update
	_, cmd := m.spinner.Update(m.spinner.Tick())
	if cmd == nil {
		t.Error("Spinner tick should return a command")
	}
}

// TestNavigationFlowDirect tests navigation by directly calling Update methods
func TestNavigationFlowDirect(t *testing.T) {
	m := createTestModelWithData(t)

	// Start in resource groups mode
	if m.core != nil && m.core.GetCurrentMode() != core.ModeResourceGroups {
		t.Errorf("Expected initial mode to be ModeResourceGroups, got %v", m.core.GetCurrentMode())
	}

	// Simulate selecting a resource group (Enter key)
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := m.Update(enterMsg)
	m = updatedModel.(model)

	var currentMode core.Mode
	if m.core != nil {
		currentMode = m.core.GetCurrentMode()
	}
	t.Logf("After Enter in resource groups mode: mode=%v, cmd=%v", currentMode, cmd != nil)

	// Test ESC key (back navigation)
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd = m.Update(escMsg)
	m = updatedModel.(model)

	if m.core != nil {
		currentMode = m.core.GetCurrentMode()
	}
	t.Logf("After ESC: mode=%v, cmd=%v", currentMode, cmd != nil)
}

// TestTeatestNavigation tests navigation using teatest
func TestTeatestNavigation(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 40))
	defer func() {
		if err := tm.Quit(); err != nil {
			t.Logf("Failed to quit test model: %v", err)
		}
	}()

	// Send window size to initialize properly
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 40})
	time.Sleep(50 * time.Millisecond)

	// Test basic key navigation
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(50 * time.Millisecond)

	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	time.Sleep(50 * time.Millisecond)

	// Test arrow keys for table navigation
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(50 * time.Millisecond)

	tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	time.Sleep(50 * time.Millisecond)

	t.Log("Navigation test completed successfully")
}

// TestResourceGroupsMode tests the resource groups mode specifically
func TestResourceGroupsMode(t *testing.T) {
	m := createTestModelWithData(t)

	// Verify we start in resource groups mode
	if m.core != nil && m.core.GetCurrentMode() != core.ModeResourceGroups {
		t.Errorf("Expected mode to be ModeResourceGroups, got %v", m.core.GetCurrentMode())
	}

	// Test refresh key
	refreshMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("r"),
	}
	updatedModel, _ := m.Update(refreshMsg)
	m = updatedModel.(model)

	// Test filter key
	filterMsg := tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("/"),
	}
	updatedModel, _ = m.Update(filterMsg)
	m = updatedModel.(model)

	t.Logf("Filter key handled")
}
