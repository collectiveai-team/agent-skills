package main

import (
	"fmt"
	"log"
	"os"

	"installer/internal/ui"
	"installer/internal/updater"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	checkForUpdate()

	program := tea.NewProgram(ui.NewModel())
	if _, err := program.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}

// checkForUpdate checks GitHub for a newer release and self-updates the binary.
// Any error or "no update" case silently returns so the TUI can start.
// Output goes to stderr since the TUI hasn't started yet.
func checkForUpdate() {
	result, err := updater.CheckForUpdate(version)
	if err != nil || !result.Available {
		return
	}

	fmt.Fprintf(os.Stderr, "  Update available: %s → %s\n", result.CurrentVer, result.LatestVer)

	if err := updater.DownloadAndReplace(result.Asset); err != nil {
		fmt.Fprintf(os.Stderr, "  Update failed: %v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "  Updated successfully. Please re-run agent-setup.\n")
	os.Exit(0)
}
