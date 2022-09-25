package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	NOTEWIDTH  int = 24
	NOTEHEIGHT int = 12
)

func main() {
	defaultColors := []lipgloss.Color{lipgloss.Color("202")}
	p := tea.NewProgram(initialNote(defaultColors, 0))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
