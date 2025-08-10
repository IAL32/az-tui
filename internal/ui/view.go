package ui

import "github.com/charmbracelet/lipgloss"

// ---------------------------- View ------------------------------

func (m model) View() string {
	if m.loading {
		return styleTitle.Render("Loading apps… ") + m.spin.View()
	}
	if m.err != nil {
		return StyleError.Render("Error: ") + m.err.Error() +
			" Press r to retry or q to quit."
	}

	var left string
	if m.mode == modeRevs {
		left = m.revList.View()
	} else {
		left = m.list.View()
	}

	// Titles
	detailsTitle := styleTitle.Render("Details")
	revsTitle := styleTitle.Render("Revisions")

	right := lipgloss.JoinVertical(
		lipgloss.Left,
		detailsTitle+"\n"+m.jsonView.View(),
		revsTitle+"\n"+m.revTable.View(),
	)
	var help string
	if m.mode == modeRevs {
		help = styleAccent.Render("[enter/e] exec  [l] logs  [b/esc] back  [q] quit  (/ filter)")
	} else {
		help = styleAccent.Render("[enter] revisions  [l] logs  [s] exec  [r] refresh  [R] reload revs  [q] quit  (/ filter)")
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(34).Render(left),
		lipgloss.NewStyle().Padding(0, 1).Render(right),
	) + "\n" + help + "\n"
}
