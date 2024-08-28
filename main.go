package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	logFile := configureLogging()
	defer logFile.Close()

	tui := tea.NewProgram(initialModel(), tea.WithAltScreen())
	_, err := tui.Run()

	if err != nil {
		log.Printf("a user interface error occurred: %s", err)
	}
}
