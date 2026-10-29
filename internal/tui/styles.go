package tui

import "github.com/charmbracelet/lipgloss"

type styles struct {
	title         lipgloss.Style
	date          lipgloss.Style
	prayer        lipgloss.Style
	next          lipgloss.Style
	nextDuringNow lipgloss.Style
	footer        lipgloss.Style
	error         lipgloss.Style
}

func defaultStyles() styles {
	return styles{
		title:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F4C95D")),
		date:          lipgloss.NewStyle().Foreground(lipgloss.Color("#8FA3B8")),
		prayer:        lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#536273")).Padding(1, 2).Width(12).Align(lipgloss.Center),
		next:          lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#111827")).Background(lipgloss.Color("#A7E8BD")).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#A7E8BD")).Padding(1, 2).Width(12).Align(lipgloss.Center),
		nextDuringNow: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Background(lipgloss.Color("#2878C8")).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#2878C8")).Padding(1, 2).Width(12).Align(lipgloss.Center),
		footer:        lipgloss.NewStyle().Foreground(lipgloss.Color("#8FA3B8")),
		error:         lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6B6B")),
	}
}
