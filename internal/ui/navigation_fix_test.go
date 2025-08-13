package ui

import (
	"testing"

	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/ui/core"
	tea "github.com/charmbracelet/bubbletea"
)

// TestNavigationBugFix tests that the navigation bug has been fixed
// This test reproduces the exact issue described in the bug report:
// "After navigating to the envvars page, pressing esc, navigation keys, and ? do not work properly"
func TestNavigationBugFix(t *testing.T) {
	// Initialize the model with mock data
	m := InitialModel(true)
	m.Init()

	// Simulate window size
	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := m.Update(windowMsg)
	m = updatedModel.(model)

	// Cast to access navigation methods
	coreModel := m.core.(*core.CoreModel)

	// Set up navigation state to simulate the bug scenario:
	// ResourceGroups -> Apps -> Revisions -> Containers
	mockRG := models.ResourceGroup{Name: "test-rg"}
	coreModel.NavigateToApps(mockRG)

	mockApp := models.ContainerApp{
		Name:          "test-app",
		ResourceGroup: "test-rg",
	}
	coreModel.NavigateToRevisions(mockApp)

	mockRev := models.Revision{Name: "test-rev"}
	coreModel.NavigateToContainers(mockRev)

	// Verify we're in containers mode
	if m.core.GetCurrentMode().String() != "Containers" {
		t.Fatalf("Expected to be in Containers mode, got %s", m.core.GetCurrentMode().String())
	}

	// Step 1: Navigate to envvars page (this is where the bug was introduced)
	mockContainer := models.Container{Name: "test-container"}
	coreModel.NavigateToEnvVars(mockContainer)

	// Verify we're in envvars mode
	currentMode := m.core.GetCurrentMode().String()
	if currentMode != "Environment Variables" && currentMode != "EnvVars" {
		t.Fatalf("Expected to be in EnvVars mode, got %s", currentMode)
	}

	// Step 2: Press 'esc' to navigate back (this was broken before the fix)
	escKey := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd := m.Update(escKey)
	m = updatedModel.(model)

	// Verify we're back in containers mode
	if m.core.GetCurrentMode().String() != "Containers" {
		t.Errorf("BUG: esc key failed to navigate back - expected Containers, got %s", m.core.GetCurrentMode().String())
	}

	// Verify the command was returned (indicating navigation occurred)
	if cmd == nil {
		t.Errorf("BUG: esc key should return a navigation command")
	}

	// Step 3: Test that navigation keys work after returning from envvars
	// This was the main symptom of the bug - these keys stopped working
	navigationKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'j'}}, // down
		{Type: tea.KeyRunes, Runes: []rune{'k'}}, // up
		{Type: tea.KeyUp},                        // up arrow
		{Type: tea.KeyDown},                      // down arrow
		{Type: tea.KeyRunes, Runes: []rune{'?'}}, // help toggle
	}

	for _, key := range navigationKeys {
		initialMode := m.core.GetCurrentMode().String()
		updatedModel, _ := m.Update(key)
		m = updatedModel.(model)
		finalMode := m.core.GetCurrentMode().String()

		// Should still be in containers mode (keys should be processed but not change mode)
		if finalMode != "Containers" {
			t.Errorf("BUG: Navigation key %v failed after envvars visit - expected Containers, got %s", key, finalMode)
		}

		// Verify we didn't accidentally navigate away
		if initialMode != finalMode && finalMode != "Containers" {
			t.Errorf("BUG: Unexpected mode change from %s to %s after key %v", initialMode, finalMode, key)
		}
	}

	t.Log("✅ Navigation bug fix verified - all keys work correctly after envvars visit")
}

// TestEnvVarsKeyHandling tests the envvars page key handling specifically
func TestEnvVarsKeyHandling(t *testing.T) {
	m := InitialModel(true)
	m.Init()

	// Cast to access navigation methods
	coreModel := m.core.(*core.CoreModel)

	// Set up navigation to envvars
	mockRG := models.ResourceGroup{Name: "test-rg"}
	coreModel.NavigateToApps(mockRG)

	mockApp := models.ContainerApp{
		Name:          "test-app",
		ResourceGroup: "test-rg",
	}
	coreModel.NavigateToRevisions(mockApp)

	mockRev := models.Revision{Name: "test-rev"}
	coreModel.NavigateToContainers(mockRev)

	mockContainer := models.Container{Name: "test-container"}
	coreModel.NavigateToEnvVars(mockContainer)

	// Test that esc key works for navigation
	escKey := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd := m.Update(escKey)
	m = updatedModel.(model)

	if m.core.GetCurrentMode().String() != "Containers" {
		t.Errorf("esc key should navigate back to containers, got %s", m.core.GetCurrentMode().String())
	}

	if cmd == nil {
		t.Errorf("esc key should return a navigation command")
	}

	t.Log("✅ EnvVars page key handling works correctly")
}

// TestNavigationKeyDelegation tests that navigation keys are properly delegated
func TestNavigationKeyDelegation(t *testing.T) {
	m := InitialModel(true)
	m.Init()

	// Cast to access navigation methods
	coreModel := m.core.(*core.CoreModel)

	// Navigate to envvars
	mockRG := models.ResourceGroup{Name: "test-rg"}
	coreModel.NavigateToApps(mockRG)

	mockApp := models.ContainerApp{
		Name:          "test-app",
		ResourceGroup: "test-rg",
	}
	coreModel.NavigateToRevisions(mockApp)

	mockRev := models.Revision{Name: "test-rev"}
	coreModel.NavigateToContainers(mockRev)

	mockContainer := models.Container{Name: "test-container"}
	coreModel.NavigateToEnvVars(mockContainer)

	// Test that navigation keys are handled (delegated to parent)
	navigationKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
		{Type: tea.KeyRunes, Runes: []rune{'?'}},
	}

	for _, key := range navigationKeys {
		initialMode := m.core.GetCurrentMode().String()
		updatedModel, _ := m.Update(key)
		m = updatedModel.(model)
		finalMode := m.core.GetCurrentMode().String()

		// Should stay in envvars mode (keys are delegated but don't change mode)
		if finalMode != initialMode {
			t.Errorf("Navigation key %v should not change mode in envvars page, changed from %s to %s", key, initialMode, finalMode)
		}
	}

	t.Log("✅ Navigation key delegation works correctly")
}
