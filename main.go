package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"time"
)

type noopWriter struct {
}

func (w noopWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func main() {
	tui := tea.NewProgram(initialModel(), tea.WithAltScreen())

	log.SetFlags(log.LstdFlags | log.LUTC)
	log.SetOutput(noopWriter{})

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for time := range ticker.C {
			tui.Send(updateTimeMsg{time: time})
		}
	}()

	_, err := tui.Run()

	if err != nil {
		log.Printf("a user interface error occurred: %s", err)
	}
}
