package ui

func (m model) View() string {
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
