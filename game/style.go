package game

import . "github.com/charmbracelet/lipgloss"

type colorFunc func(s string) string

func fg(color string) colorFunc {
	return NewStyle().Foreground(Color(color)).Render
}

var Cyan = fg("6")
var Faint = fg("8")
var Magenta = fg("5")
var Red = fg("1")
var Help = fg("241")
