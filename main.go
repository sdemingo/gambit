package main

import (
	"fmt"
	"log"
	"os"

	"gambit/game"
	"gambit/netcon"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	err := netcon.ConnectToServer("sdemingo")
	if err != nil {
		fmt.Println("Error: Server not found")
		os.Exit(1)
	}

	/*match := ""
	if len(os.Args) > 1 {
		match = os.Args[1]
	}*/
	matchName := netcon.CreateMatch()
	if matchName == "" {
		fmt.Println("Error: New match cannot be created in the server")
		os.Exit(1)
	}

	/*
		No arrancar la interfaz hasta que no tengamos todos dos jugadores

	*/

	p := tea.NewProgram(
		game.InitialModel(matchName, "sdemingo", ""),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	err = p.Start()
	if err != nil {
		log.Fatal(err)
	}
}
