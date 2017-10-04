package orgnetsim

import "reflect"

//An AgentWithMemory is a node in the network that has memory
type AgentWithMemory struct {
	AgentState
	Memory map[Color]struct{} `json:"-"`
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

//SetColor changes the color of the current Agent and counts the number of times the Agent changes color
//It also adds each color to a memory so that once it changes it's mind it doesn't change back
func (a *AgentWithMemory) SetColor(color Color) {
	_, rem := a.Memory[color]
	if !rem && a.Color != color {
		a.ChangeCount++
		a.Color = color
		a.Memory[color] = struct{}{}
	}
}

//Initialise ensures the agent is correctly initialised
func (a *AgentWithMemory) Initialise() {
	a.AgentState.Initialise()
	a.Memory = make(map[Color]struct{}, MaxColors)
	a.Memory[Grey] = struct{}{}
	a.Type = reflect.TypeOf(a).Elem().Name()
}
