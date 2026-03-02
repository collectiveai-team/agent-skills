package main

import (
	"log"

	"installer/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	program := tea.NewProgram(ui.NewModel())
	if _, err := program.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
