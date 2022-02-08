package main

import (
	"log"

	"gambit/game"
	"gambit/netcon"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	netcon.ConnectToServer("sdemingo")
	match := netcon.JoinMatch("")

	p := tea.NewProgram(
		game.InitialModel(match),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	err := p.Start()
	if err != nil {
		log.Fatal(err)
	}
}
