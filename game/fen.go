package game

import (
	"fmt"
	"strings"
)

// Display maps pieces from their FEN representations to their ASCII
// representations for a more human readable experience.
var Display = map[string]string{
	"":  " ",
	"B": "♗",
	"K": "♔",
	"N": "♘",
	"P": "♙",
	"Q": "♕",
	"R": "♖",
	"b": "♝",
	"k": "♚",
	"n": "♞",
	"p": "♟",
	"q": "♛",
	"r": "♜",
}

var DisplayASCII = map[string]string{
	"":  " ",
	"B": "A",
	"K": "R",
	"N": "C",
	"P": "P",
	"Q": "D",
	"R": "T",
	"b": "a",
	"k": "r",
	"n": "c",
	"p": "p",
	"q": "d",
	"r": "t",
}

// Tokens returns the (6) tokens of a FEN string
//
// [Pieces, Turn, Castling, En passant, Halfmove Clock, Fullmove number]
func Tokens(fen string) []string {
	return strings.Split(fen, " ")
}

// Ranks returns a slice of ranks from the first token of a FEN string
func Ranks(fen string) []string {
	return strings.Split(Tokens(fen)[0], "/")
}

// Grid returns a 8x8 grid of the board represented by the FEN string
func Grid(fen string) [8][8]string {
	var grid [8][8]string
	for r, rank := range Ranks(fen) {
		var row [8]string
		c := 0
		for _, col := range rank {
			skip := 1
			if isNumeric(col) {
				skip = runeToInt(col)
			} else {
				row[c] = fmt.Sprintf("%c", col)
			}
			c += skip
		}
		grid[r] = row
	}
	return grid
}

// isNumeric returns whether a rune is a number
func isNumeric(r rune) bool {
	return r >= '0' && r <= '9'
}

// runeToInt converts a rune to an integer
func runeToInt(r rune) int {
	return int(r - '0')
}
