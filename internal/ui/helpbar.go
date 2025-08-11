package ui

import "github.com/charmbracelet/bubbles/key"

// Help bar component factory - creates mode-specific help on the fly
func (m model) createHelpBar() string {
	// Check if we're in context selection mode
	if m.showContextList {
		return m.help.View(m.getContextHelpKeys())
	}

	// Get mode-specific key bindings from respective page files
	var modeKeys keyMap

	switch m.mode {
	case modeApps:
		modeKeys = m.getAppsHelpKeys()
	case modeRevs:
		modeKeys = m.getRevisionsHelpKeys()
	case modeContainers:
		modeKeys = m.getContainersHelpKeys()
	case modeEnvVars:
		modeKeys = m.getEnvVarsHelpKeys()
	case modeResourceGroups:
		modeKeys = m.getResourceGroupsHelpKeys()
	default:
		modeKeys = m.keys
	}

	return m.help.View(modeKeys)
}

// getContextHelpKeys returns key bindings for context selection mode
func (m model) getContextHelpKeys() contextHelpKeyMap {
	return contextHelpKeyMap{
		Enter:     m.keys.Enter,
		Back:      m.keys.Back,
		UpCombo:   m.keys.UpCombo,
		DownCombo: m.keys.DownCombo,
		Help:      m.keys.Help,
		Quit:      m.keys.Quit,
	}
}

// contextHelpKeyMap is a specialized keymap for context selection mode
type contextHelpKeyMap struct {
	Enter     key.Binding
	Back      key.Binding
	UpCombo   key.Binding
	DownCombo key.Binding
	Help      key.Binding
	Quit      key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view for context mode
func (k contextHelpKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

// FullHelp returns keybindings for the expanded help view for context mode
func (k contextHelpKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.UpCombo, k.DownCombo, k.Enter, k.Back},
		{k.Help, k.Quit},
	}
}
