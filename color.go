package orgnetsim

import "math/rand"

//A Color of an Agent
type Color int

//A List of named colours for Agents
const (
	Grey Color = iota
	Red
	Green
	Blue
	Yellow
	Orange
	Purple
)

//MaxColors is the number of defined colors
const MaxColors int = 7

//go:generate stringer -type=Color

//RandomlySelectAlternateColor selects a Color other than the one passed and other than Grey
func RandomlySelectAlternateColor(color Color) Color {
	altColor := Color(rand.Intn(MaxColors))
	for altColor == color || altColor == 0 {
		altColor = Color(rand.Intn(MaxColors))
	}
	return altColor
}
