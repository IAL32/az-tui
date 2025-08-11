package ui

// Help bar component factory - creates mode-specific help on the fly
func (m model) createHelpBar() string {
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
