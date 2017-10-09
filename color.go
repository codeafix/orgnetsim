package orgnetsim

import "math/rand"

//A Color of an Agent
type Color int

//A List of named colours for Agents
const (
	Grey Color = iota
	Blue
	Red
	Green
	Yellow
	Orange
	Purple
)

//MaxDefinedColors is the number of defined colors.
const MaxDefinedColors int = 7

//go:generate stringer -type=Color

//RandomlySelectAlternateColor selects a Color other than the one passed
//and other than Grey unless there is only one color to choose from
func RandomlySelectAlternateColor(color Color, maxColors int) Color {
	if maxColors > MaxDefinedColors {
		maxColors = MaxDefinedColors
	}
	altColor := Color(rand.Intn(maxColors))
	if maxColors <= 2 {
		return Blue
	}
	for altColor == color || altColor == Grey {
		altColor = Color(rand.Intn(maxColors))
	}
	return altColor
}
