package ui

// Help bar component factory - creates mode-specific help on the fly
func (m model) createHelpBar() string {
	// Create mode-specific key bindings on the fly
	var modeKeys keyMap

	switch m.mode {
	case modeApps:
		modeKeys = keyMap{
			Enter:       m.keys.Enter,
			Refresh:     m.keys.Refresh,
			Filter:      m.keys.Filter,
			ScrollLeft:  m.keys.ScrollLeft,
			ScrollRight: m.keys.ScrollRight,
			Help:        m.keys.Help,
			Quit:        m.keys.Quit,
		}
	case modeRevs:
		modeKeys = keyMap{
			Enter:       m.keys.Enter,
			Refresh:     m.keys.Refresh,
			RestartRev:  m.keys.RestartRev,
			Filter:      m.keys.Filter,
			ScrollLeft:  m.keys.ScrollLeft,
			ScrollRight: m.keys.ScrollRight,
			Help:        m.keys.Help,
			Back:        m.keys.Back,
			Quit:        m.keys.Quit,
		}
	case modeContainers:
		modeKeys = keyMap{
			Enter:       m.keys.Enter,
			Refresh:     m.keys.Refresh,
			Filter:      m.keys.Filter,
			Logs:        m.keys.Logs,
			Exec:        m.keys.Exec,
			ScrollLeft:  m.keys.ScrollLeft,
			ScrollRight: m.keys.ScrollRight,
			Help:        m.keys.Help,
			Back:        m.keys.Back,
			Quit:        m.keys.Quit,
		}
	case modeEnvVars:
		// Environment variables mode is read-only - no action keys
		modeKeys = keyMap{
			Filter:      m.keys.Filter,
			ScrollLeft:  m.keys.ScrollLeft,
			ScrollRight: m.keys.ScrollRight,
			Help:        m.keys.Help,
			Back:        m.keys.Back,
			Quit:        m.keys.Quit,
		}
	default:
		modeKeys = m.keys
	}

	return m.help.View(modeKeys)
}
