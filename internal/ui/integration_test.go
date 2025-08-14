package ui

import (
	"bytes"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
)

// Test that the app starts and displays resource groups
func TestAppStartup(t *testing.T) {
	m := InitialModel(true) // Use mock mode
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 30))
	defer tm.Quit()

	// Send window size to initialize properly
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 30})

	// Wait for the app to load and display content
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			// Look for typical startup content
			return bytes.Contains(bts, []byte("Resource Groups")) ||
				bytes.Contains(bts, []byte("rg-")) ||
				bytes.Contains(bts, []byte("Name")) ||
				bytes.Contains(bts, []byte("Location"))
		},
		teatest.WithCheckInterval(time.Millisecond*100),
		teatest.WithDuration(time.Second*5),
	)
}

// Test basic navigation flow
func TestBasicNavigation(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 30))
	defer tm.Quit()

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 30})

	// Wait for initial load
	time.Sleep(200 * time.Millisecond)

	// Try navigation with Enter key
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})
	time.Sleep(100 * time.Millisecond)

	// Try back navigation with ESC
	tm.Send(tea.KeyMsg{Type: tea.KeyEsc})
	time.Sleep(100 * time.Millisecond)

	// Test passes if no panic occurs
}

// Test filtering functionality
func TestFiltering(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 30))
	defer tm.Quit()

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 30})
	time.Sleep(200 * time.Millisecond)

	// Start filtering
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("/"),
	})

	// Type filter text
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("prod"),
	})

	// Apply filter
	tm.Send(tea.KeyMsg{Type: tea.KeyEnter})

	// Wait for filter to be applied
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			// Look for filter being applied or filter input
			return bytes.Contains(bts, []byte("prod")) ||
				bytes.Contains(bts, []byte("Filter"))
		},
		teatest.WithCheckInterval(time.Millisecond*50),
		teatest.WithDuration(time.Second*2),
	)
}

// Test help display
func TestHelpDisplay(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 30))
	defer tm.Quit()

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 30})
	time.Sleep(200 * time.Millisecond)

	// Toggle help
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("?"),
	})

	// Wait for help to appear
	teatest.WaitFor(
		t, tm.Output(),
		func(bts []byte) bool {
			return bytes.Contains(bts, []byte("help")) ||
				bytes.Contains(bts, []byte("Help")) ||
				bytes.Contains(bts, []byte("enter")) ||
				bytes.Contains(bts, []byte("quit")) ||
				bytes.Contains(bts, []byte("Key"))
		},
		teatest.WithCheckInterval(time.Millisecond*50),
		teatest.WithDuration(time.Second*2),
	)
}

// Test refresh functionality
func TestRefresh(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 30))
	defer tm.Quit()

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 30})
	time.Sleep(200 * time.Millisecond)

	// Send refresh command
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("r"),
	})

	// Wait a bit to see if refresh works without errors
	time.Sleep(200 * time.Millisecond)

	// Test passes if no panic occurs
}

// Test window resizing
func TestWindowResize(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))
	defer tm.Quit()

	// Initialize with small size
	tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})
	time.Sleep(100 * time.Millisecond)

	// Resize to larger
	tm.Send(tea.WindowSizeMsg{Width: 150, Height: 40})
	time.Sleep(100 * time.Millisecond)

	// Resize to smaller
	tm.Send(tea.WindowSizeMsg{Width: 60, Height: 20})
	time.Sleep(100 * time.Millisecond)

	// Test passes if no panic occurs during resizing
}

// Test quit functionality
func TestQuit(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})
	time.Sleep(100 * time.Millisecond)

	// Send quit command
	tm.Send(tea.KeyMsg{
		Type:  tea.KeyRunes,
		Runes: []rune("q"),
	})

	// Wait for the model to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}

// Test Ctrl+C quit
func TestCtrlCQuit(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(80, 24))

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 80, Height: 24})
	time.Sleep(100 * time.Millisecond)

	// Send Ctrl+C
	tm.Send(tea.KeyMsg{Type: tea.KeyCtrlC})

	// Wait for the model to finish
	tm.WaitFinished(t, teatest.WithFinalTimeout(time.Second*2))
}

// Test arrow key navigation
func TestArrowKeyNavigation(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 30))
	defer tm.Quit()

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 30})
	time.Sleep(200 * time.Millisecond)

	// Test arrow keys
	tm.Send(tea.KeyMsg{Type: tea.KeyDown})
	time.Sleep(50 * time.Millisecond)

	tm.Send(tea.KeyMsg{Type: tea.KeyUp})
	time.Sleep(50 * time.Millisecond)

	tm.Send(tea.KeyMsg{Type: tea.KeyLeft})
	time.Sleep(50 * time.Millisecond)

	tm.Send(tea.KeyMsg{Type: tea.KeyRight})
	time.Sleep(50 * time.Millisecond)

	// Test passes if no panic occurs
}

// Test that the app handles rapid key presses without crashing
func TestRapidKeyPresses(t *testing.T) {
	m := InitialModel(true)
	tm := teatest.NewTestModel(t, m, teatest.WithInitialTermSize(120, 30))
	defer tm.Quit()

	// Initialize
	tm.Send(tea.WindowSizeMsg{Width: 120, Height: 30})
	time.Sleep(100 * time.Millisecond)

	// Send rapid key presses
	keys := []tea.KeyMsg{
		{Type: tea.KeyDown},
		{Type: tea.KeyUp},
		{Type: tea.KeyEnter},
		{Type: tea.KeyEsc},
		{Type: tea.KeyDown},
		{Type: tea.KeyDown},
		{Type: tea.KeyUp},
	}

	for _, key := range keys {
		tm.Send(key)
		time.Sleep(10 * time.Millisecond) // Very short delay
	}

	// Test passes if no panic occurs
}
