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

//MaxColors is the number of defined colors.
const MaxColors int = 4

//go:generate stringer -type=Color

//RandomlySelectAlternateColor selects a Color other than the one passed
//and other than Grey unless there is only one color to choose from
func RandomlySelectAlternateColor(color Color) Color {
	altColor := Color(rand.Intn(MaxColors))
	if MaxColors <= 2 {
		return Blue
	}
	for altColor == color || altColor == Grey {
		altColor = Color(rand.Intn(MaxColors))
	}
	return altColor
}
