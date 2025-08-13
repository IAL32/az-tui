package core

import (
	"github.com/IAL32/az-tui/internal/ui/layouts"
	tea "github.com/charmbracelet/bubbletea"
)

// Mode type alias to layouts.Mode for backward compatibility
type Mode = layouts.Mode

// Mode constants for backward compatibility
const (
	ModeResourceGroups = layouts.ModeResourceGroups
	ModeApps           = layouts.ModeApps
	ModeRevisions      = layouts.ModeRevisions
	ModeContainers     = layouts.ModeContainers
	ModeEnvVars        = layouts.ModeEnvVars
)

// NavigationState holds the current navigation context
type NavigationState struct {
	CurrentRG            string // Current resource group
	CurrentAppID         string // When viewing revisions
	CurrentRevName       string // When viewing containers
	CurrentContainerName string // When viewing environment variables
}

// Reset clears all navigation state
func (ns *NavigationState) Reset() {
	ns.CurrentRG = ""
	ns.CurrentAppID = ""
	ns.CurrentRevName = ""
	ns.CurrentContainerName = ""
}

// ResetFrom resets navigation state from a specific level
func (ns *NavigationState) ResetFrom(mode Mode) {
	switch mode {
	case ModeResourceGroups:
		ns.Reset()
	case ModeApps:
		ns.CurrentAppID = ""
		ns.CurrentRevName = ""
		ns.CurrentContainerName = ""
	case ModeRevisions:
		ns.CurrentRevName = ""
		ns.CurrentContainerName = ""
	case ModeContainers:
		ns.CurrentContainerName = ""
	}
}

// GetBreadcrumb returns a breadcrumb string for the current navigation state
func (ns *NavigationState) GetBreadcrumb() string {
	var parts []string

	if ns.CurrentRG != "" {
		parts = append(parts, ns.CurrentRG)
	}
	if ns.CurrentAppID != "" {
		parts = append(parts, ns.CurrentAppID)
	}
	if ns.CurrentRevName != "" {
		parts = append(parts, ns.CurrentRevName)
	}
	if ns.CurrentContainerName != "" {
		parts = append(parts, ns.CurrentContainerName)
	}

	if len(parts) == 0 {
		return ""
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += " > " + parts[i]
	}
	return result
}

// CoreEvent represents events that the core system can handle
type CoreEvent interface {
	tea.Msg
}

// NavigationEvent represents navigation-related events
type NavigationEvent struct {
	FromMode Mode
	ToMode   Mode
	Data     interface{}
}

// ModeChangeEvent represents mode change events
type ModeChangeEvent struct {
	OldMode Mode
	NewMode Mode
}

// PageLoadEvent represents page loading events
type PageLoadEvent struct {
	Mode    Mode
	Loading bool
	Error   error
}

// DataUpdateEvent represents data update events
type DataUpdateEvent struct {
	Mode Mode
	Data interface{}
}
