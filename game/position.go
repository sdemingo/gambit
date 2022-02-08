package game

import (
	"fmt"
	"strconv"
)

// colToFile returns the file given a column
func colToFile(col int) string {
	if col < FirstCol {
		col = FirstCol
	} else if col > LastCol {
		col = LastCol
	}
	return fmt.Sprintf("%c", col+'a')
}

// rowToRank returns a rank given a row
func rowToRank(row int) int {
	if row < FirstRow {
		row = FirstRow
	} else if row > LastRow {
		row = LastRow
	}
	return row + 1
}

// ToSquare returns the square position (e.g. a1) of a given row and column
// (e.g. 0,0) for display or checking legal moves.
func ToSquare(row, col int) string {
	return colToFile(col) + strconv.Itoa(rowToRank(row))
}
