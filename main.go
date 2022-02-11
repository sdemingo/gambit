package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"

	"gambit/game"
	"gambit/netcon"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {

	ascciflag := flag.Bool("a", false, "No usar símbolos Unicode para las piezas")
	helpflag := flag.Bool("h", false, "Mostrar este mensaje de ayuda")
	userflag := flag.String("user", "", "user name in the game")
	gameflag := flag.String("game", "", "id of the game")

	flag.Parse()

	if *helpflag {
		flag.PrintDefaults()
		os.Exit(1)
	}

	username := ""
	if *userflag == "" {
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		username = user.Username
	} else {
		username = *userflag
	}

	netcon.InitLog(username)

	err := netcon.ConnectToServer(username)
	if err != nil {
		fmt.Println("Error: Server not found")
		os.Exit(1)
	}

	playerBlack := ""
	playerWhite := ""
	matchName := ""

	if *gameflag == "" {
		/*
			Create a new game and play with whites
		*/
		playerWhite = username
		matchName, err = netcon.CreateMatch()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("\n\nBienvenido %s. Has creado una partida en el servidor.\n\n", playerWhite)
		fmt.Printf("Su identificador es: %s\n", matchName)
		fmt.Println("Comunicáselo a tu oponente y espera a que se conecte .... ")

		playerBlack, err = netcon.ReceivePlayerName()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else {
		/*
			Join to a created game and play with black
		*/
		playerBlack = username
		matchName = *gameflag
		playerWhite, err = netcon.JoinMatch(matchName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	/*
		Run UI

	*/

	p := tea.NewProgram(
		game.InitialModel(matchName, playerWhite, playerBlack, (playerWhite == username), *ascciflag),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	err = p.Start()
	if err != nil {
		log.Fatal(err)
	}
}
