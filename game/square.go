package game

import (
	"strconv"
)

// fileToCol returns column number (e.g. 0) for a given file (e.g. 'a').
func fileToCol(file rune) int {
	col := int(file - 'a')
	if col < FirstCol {
		col = FirstCol
	} else if col > LastCol {
		col = LastCol
	}
	return col
}

// rankToRow returns a row number (e.g. 0) for a given rank (e.g. 1).
func rankToRow(rank int) int {
	row := rank - 1
	if row < FirstRow {
		row = FirstRow
	} else if row > LastRow {
		row = LastRow
	}
	return row
}

// ToPosition takes a square (e.g. a1) and returns the corresponding row and
// column (e.g. 0,0) for compatibility with the grid (8x8 matrix).
func ToPosition(square string) (int, int) {
	col := fileToCol(rune(square[0]))
	row, _ := strconv.Atoi(string(square[1]))
	row = rankToRow(row)
	return col, row
}
