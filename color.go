package orgnetsim

//A Color of an Agent
type Color int

//A List of named colours for Agents
const (
	Grey Color = iota
	Black
	White
	Red
	Green
	Blue
	Yellow
	Orange
	Purple
)

//go:generate stringer -type=Color
