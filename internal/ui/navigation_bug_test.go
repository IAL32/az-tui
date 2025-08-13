package ui

import (
	"testing"

	"github.com/IAL32/az-tui/internal/models"
	"github.com/IAL32/az-tui/internal/ui/core"
	tea "github.com/charmbracelet/bubbletea"
)

// TestEnvVarsPageKeyDelegation tests that the envvars page properly delegates
// navigation keys to the parent system
func TestEnvVarsPageKeyDelegation(t *testing.T) {
	// Initialize the model with mock data
	m := InitialModel(true)

	// Initialize the model
	m.Init()

	// Simulate window size
	windowMsg := tea.WindowSizeMsg{Width: 100, Height: 30}
	updatedModel, _ := m.Update(windowMsg)
	m = updatedModel.(model)

	// We need to access the core model directly to test navigation
	// Since the CoreInterface doesn't expose navigation methods, we'll cast it
	coreModel := m.core.(*core.CoreModel)

	// Step 1: Set up navigation state manually to simulate being in containers mode
	// This simulates the navigation chain: ResourceGroups -> Apps -> Revisions -> Containers
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

	// Step 2: Test navigation keys work on containers page BEFORE visiting envvars
	t.Log("Testing navigation keys on containers page BEFORE envvars visit")

	testKeys := []tea.KeyMsg{
		{Type: tea.KeyRunes, Runes: []rune{'j'}},
		{Type: tea.KeyRunes, Runes: []rune{'k'}},
		{Type: tea.KeyRunes, Runes: []rune{'?'}},
		{Type: tea.KeyUp},
		{Type: tea.KeyDown},
	}

	for _, keyMsg := range testKeys {
		updatedModel, _ := m.Update(keyMsg)
		m = updatedModel.(model)

		// Should still be in containers mode
		if m.core.GetCurrentMode().String() != "Containers" {
			t.Errorf("Expected to stay in Containers mode after key %v, got %s", keyMsg, m.core.GetCurrentMode().String())
		}
	}

	// Step 3: Navigate to envvars page
	t.Log("Navigating to envvars page")

	mockContainer := models.Container{Name: "test-container"}
	coreModel.NavigateToEnvVars(mockContainer)

	// Verify we're in envvars mode
	if m.core.GetCurrentMode().String() != "EnvVars" {
		t.Fatalf("Expected to be in EnvVars mode, got %s", m.core.GetCurrentMode().String())
	}

	// Step 4: Test key handling on envvars page
	t.Log("Testing key handling on envvars page")

	// Test that navigation keys are handled (should stay in envvars but be processed)
	for _, keyMsg := range testKeys {
		if keyMsg.Type == tea.KeyRunes && string(keyMsg.Runes) == "?" {
			continue // Skip help key for now
		}

		updatedModel, _ := m.Update(keyMsg)
		m = updatedModel.(model)

		// Should still be in envvars mode (keys should be processed but not change mode)
		if m.core.GetCurrentMode().String() != "EnvVars" {
			t.Errorf("Expected to stay in EnvVars mode after key %v, got %s", keyMsg, m.core.GetCurrentMode().String())
		}
	}

	// Step 5: Navigate back from envvars using 'esc'
	t.Log("Navigating back from envvars using 'esc'")

	escKey := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ = m.Update(escKey)
	m = updatedModel.(model)

	// Should be back in containers mode
	if m.core.GetCurrentMode().String() != "Containers" {
		t.Fatalf("Expected to be back in Containers mode after esc, got %s", m.core.GetCurrentMode().String())
	}

	// Step 6: Test navigation keys AFTER visiting envvars - THIS IS WHERE THE BUG OCCURS
	t.Log("Testing navigation keys on containers page AFTER envvars visit - this should reveal the bug")

	for _, keyMsg := range testKeys {
		t.Logf("Testing key: %v", keyMsg)

		initialMode := m.core.GetCurrentMode().String()
		updatedModel, cmd := m.Update(keyMsg)
		m = updatedModel.(model)
		finalMode := m.core.GetCurrentMode().String()

		t.Logf("Key %v: %s -> %s, cmd: %v", keyMsg, initialMode, finalMode, cmd != nil)

		// Should still be in containers mode
		if finalMode != "Containers" {
			t.Errorf("BUG DETECTED: Expected to stay in Containers mode after key %v post-envvars, got %s", keyMsg, finalMode)
		}

		// The key should be processed (either handled by page or delegated to table)
		// We can't easily test this without more introspection, but the mode should remain stable
	}

	t.Log("Navigation test completed")
}

// TestEnvVarsPageKeyHandlingIsolated tests the envvars page key handling in isolation
func TestEnvVarsPageKeyHandlingIsolated(t *testing.T) {
	// This test focuses specifically on the envvars page key handling
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

	// Test key handling on envvars page
	testCases := []struct {
		name        string
		key         tea.KeyMsg
		expectMode  string
		description string
	}{
		{
			name:        "j_key_delegation",
			key:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			expectMode:  "EnvVars",
			description: "j key should be delegated to parent but stay in envvars",
		},
		{
			name:        "k_key_delegation",
			key:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			expectMode:  "EnvVars",
			description: "k key should be delegated to parent but stay in envvars",
		},
		{
			name:        "up_arrow_delegation",
			key:         tea.KeyMsg{Type: tea.KeyUp},
			expectMode:  "EnvVars",
			description: "up arrow should be delegated to parent but stay in envvars",
		},
		{
			name:        "down_arrow_delegation",
			key:         tea.KeyMsg{Type: tea.KeyDown},
			expectMode:  "EnvVars",
			description: "down arrow should be delegated to parent but stay in envvars",
		},
		{
			name:        "help_key_delegation",
			key:         tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}},
			expectMode:  "EnvVars",
			description: "? key should be delegated to parent but stay in envvars",
		},
		{
			name:        "esc_navigation",
			key:         tea.KeyMsg{Type: tea.KeyEsc},
			expectMode:  "Containers",
			description: "esc key should navigate back to containers",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset to envvars page if needed (except for esc test)
			if tc.name != "esc_navigation" && m.core.GetCurrentMode().String() != "EnvVars" {
				coreModel.NavigateToEnvVars(mockContainer)
			}

			initialMode := m.core.GetCurrentMode().String()
			t.Logf("Initial mode: %s", initialMode)
			t.Logf("Testing: %s", tc.description)

			// Send the key
			updatedModel, cmd := m.Update(tc.key)
			m = updatedModel.(model)

			finalMode := m.core.GetCurrentMode().String()
			t.Logf("Final mode after %s: %s", tc.name, finalMode)
			t.Logf("Command returned: %v", cmd != nil)

			if finalMode != tc.expectMode {
				t.Errorf("Expected mode %s after %s, got %s", tc.expectMode, tc.name, finalMode)
			}
		})
	}
}

// TestKeyHandlingFlowWithDiagnostics tests the complete key handling flow with diagnostic output
func TestKeyHandlingFlowWithDiagnostics(t *testing.T) {
	m := InitialModel(true)
	m.Init()

	// Cast to access navigation methods
	coreModel := m.core.(*core.CoreModel)

	// Set up the navigation chain
	mockRG := models.ResourceGroup{Name: "test-rg"}
	coreModel.NavigateToApps(mockRG)

	mockApp := models.ContainerApp{
		Name:          "test-app",
		ResourceGroup: "test-rg",
	}
	coreModel.NavigateToRevisions(mockApp)

	mockRev := models.Revision{Name: "test-rev"}
	coreModel.NavigateToContainers(mockRev)

	// Test that keys work on containers page
	t.Log("=== Testing 'j' key on containers page BEFORE envvars ===")
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, cmd := m.Update(keyMsg)
	m = updatedModel.(model)

	t.Logf("Containers page 'j' key - Mode: %s, Cmd: %v", m.core.GetCurrentMode().String(), cmd != nil)

	// Navigate to envvars
	t.Log("=== Navigating to envvars ===")
	mockContainer := models.Container{Name: "test-container"}
	coreModel.NavigateToEnvVars(mockContainer)

	// Test key on envvars page
	t.Log("=== Testing 'j' key on envvars page ===")
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, cmd = m.Update(keyMsg)
	m = updatedModel.(model)

	t.Logf("EnvVars page 'j' key - Mode: %s, Cmd: %v", m.core.GetCurrentMode().String(), cmd != nil)

	// Navigate back
	t.Log("=== Navigating back with esc ===")
	keyMsg = tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, cmd = m.Update(keyMsg)
	m = updatedModel.(model)

	t.Logf("After esc from envvars - Mode: %s, Cmd: %v", m.core.GetCurrentMode().String(), cmd != nil)

	// Test that keys work again on containers page
	t.Log("=== Testing 'j' key on containers page AFTER envvars ===")
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, cmd = m.Update(keyMsg)
	m = updatedModel.(model)

	t.Logf("Containers page 'j' key AFTER envvars - Mode: %s, Cmd: %v", m.core.GetCurrentMode().String(), cmd != nil)

	// This should work but might not due to the bug
	if m.core.GetCurrentMode().String() != "Containers" {
		t.Errorf("BUG: Navigation broken after envvars visit - expected Containers, got %s", m.core.GetCurrentMode().String())
	}

	// Test help key specifically (mentioned in bug report)
	t.Log("=== Testing '?' key on containers page AFTER envvars ===")
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}}
	updatedModel, cmd = m.Update(keyMsg)
	m = updatedModel.(model)

	t.Logf("Containers page '?' key AFTER envvars - Mode: %s, Cmd: %v", m.core.GetCurrentMode().String(), cmd != nil)

	if m.core.GetCurrentMode().String() != "Containers" {
		t.Errorf("BUG: Help key broken after envvars visit - expected Containers, got %s", m.core.GetCurrentMode().String())
	}
}
