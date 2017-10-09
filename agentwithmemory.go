package orgnetsim

import "reflect"

//An AgentWithMemory is a node in the network that has memory
type AgentWithMemory struct {
	AgentState
	PreviousColors map[Color]struct{} `json:"-"`
	ShortMemory    map[Color]struct{} `json:"-"`
	MaxColors      int                `json:"-"`
}

/*ReadMail checks for any messages it received in its own Mail queue. If it receives
one then it decides whether to update its color.
*/
func (a *AgentWithMemory) ReadMail(n RelationshipMgr) Color {
	msg, received := a.ReceiveMsg()
	if received {
		ra, isAgentWithMem := n.GetAgentByID(msg).(*AgentWithMemory)
		if isAgentWithMem {
			c, update := a.UpdateColor(n, &ra.AgentState)
			if update {
				a.SetColor(c)
			}
		}
	}
	return a.Color
}

//State returns the struct containing the state of this Agent
func (a *AgentWithMemory) State() *AgentState {
	return &a.AgentState
}

//SetColor changes the color of the current Agent and counts the number of times the Agent changes color
//It also adds each color to a short term memory. It will only change its color if it hears about another
//color twice. Once it has updated its color it will clear its short term memory and update its long term
//memory with its previous color so that it doesn't get set to the same color twice
func (a *AgentWithMemory) SetColor(color Color) {
	_, rem := a.PreviousColors[color]
	if !rem && a.Color != color {
		_, rem := a.ShortMemory[color]
		if rem {
			a.ChangeCount++
			a.PreviousColors[a.Color] = struct{}{}
			a.Color = color
			a.ShortMemory = make(map[Color]struct{}, a.MaxColors)
		} else {
			a.ShortMemory[color] = struct{}{}
		}
	}
}

//Initialise ensures the agent is correctly initialised
func (a *AgentWithMemory) Initialise(n RelationshipMgr) {
	a.AgentState.Initialise(n)
	a.MaxColors = n.MaxColors()
	a.PreviousColors = make(map[Color]struct{}, a.MaxColors)
	a.PreviousColors[Grey] = struct{}{}
	a.ShortMemory = make(map[Color]struct{}, a.MaxColors)
	a.Type = reflect.TypeOf(a).Elem().Name()
}
