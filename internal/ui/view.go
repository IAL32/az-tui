package ui

func (m model) View() string {
	// Show context list when active
	if m.showContextList {
		return m.viewContextList()
	}

	// Existing view logic remains unchanged
	switch m.mode {
	case modeContainers:
		return m.viewContainers()
	case modeRevs:
		return m.viewRevs()
	case modeEnvVars:
		return m.viewEnvVars()
	case modeResourceGroups:
		return m.viewResourceGroups()
	default:
		return m.viewApps()
	}
}
