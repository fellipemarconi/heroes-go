package main

import (
	"api/internal/admin/tui"
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	p := tea.NewProgram(tui.NewDashboard())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
