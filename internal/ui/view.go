package ui

import "github.com/charmbracelet/lipgloss"

// ---------------------------- View ------------------------------

func (m model) View() string {
	if m.loading {
		return styleTitle.Render("Loading apps… ") + m.spin.View()
	}
	if m.err != nil {
		return StyleError.Render("Error: ") + m.err.Error() + " Press r to retry or q to quit."
	}

	left := m.list.View()

	rightHeader := styleTitle.Render("Details")
	if m.activePane == paneRevisions {
		rightHeader = styleTitle.Render("Revisions")
	}

	var right string
	if m.activePane == paneDetails {
		// ensure header is present above jsonView
		right = rightHeader + " " + m.jsonView.View()
	} else {
		right = rightHeader + " " + m.revTable.View()
	}

	help := styleAccent.Render("[q] quit  [r] refresh  [tab] switch pane  [R] reload revs  [l] logs  [e] exec  (/ filter list)")

	// simple side-by-side layout
	return lipgloss.JoinHorizontal(lipgloss.Top,
		lipgloss.NewStyle().Width(34).Render(left),
		lipgloss.NewStyle().Padding(0, 1).Render(right),
	) + "\n" + help + "\n"
}
