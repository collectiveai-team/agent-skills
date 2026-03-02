package ui

import "github.com/charmbracelet/lipgloss"

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true)
}

func highlightStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("69")).Bold(true)
}

func warnStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("208")).Bold(true)
}

func helpStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
}

func progressStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
}

func indentStyle() lipgloss.Style {
	return lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("196"))
}

func dimStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242"))
}

func tabActiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("69")).Underline(true)
}

func tabInactiveStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
}

func successStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
}
