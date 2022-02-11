package game

import (
	"fmt"
	"gambit/netcon"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	dt "github.com/dylhunn/dragontoothmg"
)

const (
	Cols = 8
	Rows = 8
)

const (
	FirstCol = 0
	FirstRow = 0
	LastCol  = Cols - 1
	LastRow  = Rows - 1
)

// model stores the state of the chess game.
//
// It tracks the board, legal moves, and the selected piece. It also keeps
// track of the subset of legal moves for the currently selected piece
type model struct {
	game       string
	userWhite  string
	userBlack  string
	board      *dt.Board
	moves      []dt.Move
	pieceMoves []dt.Move
	selected   string
	buffer     string
	flipped    bool
	wait       bool // wait for the move of your opponent
	cmoves     chan string
	ascii      bool
}

// InitialModel returns an initial model of the game board. It uses the
// starting position of a normal chess game and generates the legal moves from
// the starting position.
func InitialModel(game string, white string, black string, start bool, ascii bool) tea.Model {
	board := dt.ParseFen(dt.Startpos)

	return model{
		game:      game,
		userWhite: white,
		userBlack: black,
		board:     &board,
		moves:     board.GenerateLegalMoves(),
		wait:      !start,
		cmoves:    make(chan string),
		ascii:     ascii,
	}
}

// Init Initializes the model
func (m model) Init() tea.Cmd {
	return tea.Batch(listenForMove(m.cmoves), waitForMove(m.cmoves))
}

type moveReceived string

func listenForMove(c chan string) tea.Cmd {
	return func() tea.Msg {
		for {
			moveMsg, err := netcon.ReciveMsg()
			if moveMsg != nil && moveMsg.Cmd == netcon.MOVE {
				fields := strings.Split(moveMsg.Args, ":")
				if len(fields) == 2 {
					log.Printf("Recibido movimiento: %s\n", fields[1])
					c <- fields[1]
				}
			}

			if err != nil {
				log.Println("El servidor cerró la conexión")
				c <- "FINAL:"
			}
			if moveMsg != nil && moveMsg.Cmd == netcon.END {
				log.Println("Se recibió final de la partida")
				c <- "FINAL:" + moveMsg.Args
			}
		}
	}
}

func waitForMove(c chan string) tea.Cmd {
	return func() tea.Msg {
		move := <-c
		return moveReceived(move)
	}
}

// View converts a FEN string into a human readable chess board. All pieces and
// empty squares are arranged in a grid-like pattern. The selected piece is
// highlighted and the legal moves for the selected piece are indicated by a
// dot (.) for empty squares. Pieces that may be captured by the selected piece
// are highlighted.
//
// For example, if the user selects the white pawn on E2 we indicate that they
// can move to E3 and E4 legally.
//
//    ┌───┬───┬───┬───┬───┬───┬───┬───┐
//  8 │ ♜ │ ♞ │ ♝ │ ♛ │ ♚ │ ♝ │ ♞ │ ♜ │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  7 │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │ ♟ │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  6 │   │   │   │   │   │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  5 │   │   │   │   │   │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  4 │   │   │   │   │ . │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  3 │   │   │   │   │ . │   │   │   │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  2 │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │ ♙ │
//    ├───┼───┼───┼───┼───┼───┼───┼───┤
//  1 │ ♖ │ ♘ │ ♗ │ ♕ │ ♔ │ ♗ │ ♘ │ ♖ │
//    └───┴───┴───┴───┴───┴───┴───┴───┘
//      A   B   C   D   E   F   G   H
//
func (m model) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("  Partida: %s\n", m.game))
	s.WriteString(fmt.Sprintf("  Blancas: %s\n", m.userWhite))
	s.WriteString(fmt.Sprintf("  Negras: %s\n", m.userBlack))
	if m.wait {
		s.WriteString("Esperando movimiento de tu oponente ...\n\n")
	} else {
		s.WriteString("Mueve\n\n")
	}

	s.WriteString(BorderTop())

	// Traverse through the rows and columns of the board and print out the
	// pieces and empty squares. Once a piece is selected, highlight the legal
	// moves and pieces that may be captured by the selected piece.
	var rows = Grid(m.board.ToFen())

	for r := FirstRow; r < Rows; r++ {
		row := rows[r]
		rr := LastRow - r

		if m.flipped {
			row = rows[LastRow-r]
			rr = r
		}

		s.WriteString(Faint(fmt.Sprintf(" %d ", rr+1)) + Vertical)

		for c, cell := range row {
			display := Display[cell]
			if m.ascii {
				display = cell
				if display == "" {
					display = " "
				}
			}

			// The user selected the current cell, highlight it so they know it is
			// selected.
			if m.selected == ToSquare(rr, c) {
				display = Cyan(display)
			}

			// Show all the cells to which the piece may move. If it is an empty cell
			// we present a coloured dot, otherwise color the capturable piece.
			if IsLegal(m.pieceMoves, ToSquare(rr, c)) {
				if cell == "" {
					display = "."
				}
				display = Magenta(display)
			}

			s.WriteString(fmt.Sprintf(" %s %s", display, Vertical))
		}
		s.WriteRune('\n')

		if r != LastRow {
			s.WriteString(Middle())
		}
	}

	s.WriteString(BorderBottom() + Faint(BottomLabels()))
	return s.String()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		if msg.Type != tea.MouseLeft {
			return m, nil
		}

		// Find the square the user clicked on, this will either be our square
		// square for our piece or the destination square for a move if a piece is
		// already square and that destination square completes a legal move
		square := Cell(msg.X, msg.Y, m.flipped)
		return m.Select(square)

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "ctrl+f":
			m.flipped = !m.flipped

		case "a", "b", "c", "d", "e", "f", "g", "h":
			m.buffer = msg.String()

		case "1", "2", "3", "4", "5", "6", "7", "8":
			var move string
			if m.buffer != "" {
				move = m.buffer + msg.String()
				m.buffer = ""
			}
			return m.Select(move)

		case "esc":
			return m.Deselect()
		}

	case moveReceived:
		movestr := fmt.Sprintf("%s", msg)
		if strings.HasPrefix(movestr, "FINAL:") {
			return m, tea.Quit
		}
		if m.wait {
			move, err := dt.ParseMove(movestr)
			if err == nil {
				m.board.Apply(move)
				log.Println("Movimiento aplicado")
				m.moves = m.board.GenerateLegalMoves()
			} else {
				log.Println("Error: Se intento aplicar un movimiento mal formado")
			}
		}
		m.wait = false

		return m, waitForMove(m.cmoves)
	}

	return m, nil
}

/*
func ApplyOponentMove(movestr string) tea.Msg {
	log.Println("retornamos moveReceived con " + movestr)
	return moveReceived(movestr)
}
*/
func (m model) Deselect() (tea.Model, tea.Cmd) {
	m.selected = ""
	m.pieceMoves = []dt.Move{}
	return m, nil
}

func (m model) Select(square string) (tea.Model, tea.Cmd) {
	if m.wait {
		return m, nil //it's not your turn!!
	}

	// If the user has already selected a piece, check see if the square that
	// the user clicked on is a legal move for that piece. If so, make the move.
	if m.selected != "" {
		from := m.selected
		to := square

		for _, move := range m.pieceMoves {
			if move.String() == from+to {
				m.board.Apply(move)

				// send move to server and wait for the oponent move
				netcon.SendMove(move, m.game)
				m.wait = true

				// We have applied a new move and the chess board is in a new state.
				// We must generate the new legal moves for the new state.
				m.moves = m.board.GenerateLegalMoves()

				// We have made a move, so we no longer have a selected piece or
				// legal moves for any selected pieces.
				return m.Deselect()
			}
		}

		// The user clicked on a square that wasn't a legal move for the selected
		// piece, so we select the piece that was clicked on instead
		m.selected = to
	} else {
		m.selected = square
	}

	// After a mouse click, we must generate the legal moves for the selected
	// piece, if there is a newly selected piece
	m.pieceMoves = LegalSelected(m.moves, m.selected)

	return m, nil
}
