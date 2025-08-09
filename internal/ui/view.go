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

	left := m.list.View()

	// Titles
	detailsTitle := styleTitle.Render("Details")
	revsTitle := styleTitle.Render("Revisions")

	right := lipgloss.JoinVertical(
		lipgloss.Left,
		detailsTitle+"\n"+m.jsonView.View(),
		revsTitle+"\n"+m.revTable.View(),
	)

	help := styleAccent.Render(
		"[q] quit  [r] refresh  [R] reload revs  [l] logs  [e] exec  (/ filter list)",
	)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(34).Render(left),
		lipgloss.NewStyle().Padding(0, 1).Render(right),
	) + "\n" + help + "\n"
}
