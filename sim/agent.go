package sim

import "reflect"

//An AgentState carries the state of a node in the network
type AgentState struct {
	ID             string      `json:"id"`
	Color          Color       `json:"color"`
	Susceptability float64     `json:"susceptability"`
	Influence      float64     `json:"influence"`
	Contrariness   float64     `json:"contrariness"`
	Mail           chan string `json:"-"`
	ChangeCount    int         `json:"change"`
	Type           string      `json:"type"`
	X              float64     `json:"fx,omitempty"`
	Y              float64     `json:"fy,omitempty"`
}

//Agent is an interface that allows interaction with an Agent
type Agent interface {
	Initialise(n RelationshipMgr)
	Identifier() string
	State() *AgentState
	SendMail(n RelationshipMgr) int
	ReadMail(n RelationshipMgr) Color
	ClearMail()
	GetColor() Color
	PostMsg(msg string) bool
	ReceiveMsg() (string, bool)
}

//Initialise ensures the agent is correctly initialised
func (a *AgentState) Initialise(n RelationshipMgr) {
	a.Mail = make(chan string, 1)
	a.Type = reflect.TypeOf(a).Elem().Name()
}

//Identifier returns the Identifier for the Agent
func (a *AgentState) Identifier() string {
	return a.ID
}

//State returns the struct containing the state of this Agent
func (a *AgentState) State() *AgentState {
	return a
}

/*SendMail iterates over a randomly ordered slice of related agents trying to find a match. It sends a mail to the
first successful match it finds.
*/
func (a *AgentState) SendMail(n RelationshipMgr) int {
	for _, ra := range n.GetRelatedAgents(a) {
		if ra.PostMsg(a.ID) {
			return 1
		}
	}
	return 0
}

/*ReadMail checks for any messages it received in its own Mail queue. If it receives
one then it decides whether to update its color.
*/
func (a *AgentState) ReadMail(n RelationshipMgr) Color {
	msg, received := a.ReceiveMsg()
	if received {
		ra, isAgent := n.GetAgentByID(msg).(*AgentState)
		if isAgent {
			c, update := a.UpdateColor(n, ra)
			if update {
				a.SetColor(c)
			}
		}
	}
	return a.Color
}

//UpdateColor looks at the properties of the passed agent and decides what the agent should update its color to
func (a *AgentState) UpdateColor(n RelationshipMgr, ra *AgentState) (Color, bool) {
	n.IncrementLinkStrength(a.Identifier(), ra.Identifier())
	if ra.Influence > a.Susceptability {
		if a.Contrariness > ra.Influence {
			altColor := RandomlySelectAlternateColor(a.Color, n.MaxColors())
			return altColor, true
		}
		return ra.Color, true
	}
	return Grey, false
}

//SetColor changes the color of the current Agent and counts the number of times the Agent changes color
//It also adds each color to a memory so that once it changes it's mind it doesn't change back
func (a *AgentState) SetColor(color Color) {
	if a.Color != color {
		a.ChangeCount++
		a.Color = color
	}
}

//GetColor returns the Color of this Agent
func (a *AgentState) GetColor() Color {
	return a.Color
}

//PostMsg tries to add an entry into an Agent's Mail channel, if it succeeds, that Agent will be blocked
//for any other Agent trying to send a Mail and this function returns true (the Agent is now Matched).
//If it returns false the Agent is already matched by another Agent.
func (a *AgentState) PostMsg(msg string) bool {
	select {
	case a.Mail <- msg:
		return true
	default:
		return false
	}
}

//ReceiveMsg picks a message up from the Agent's Mail channel
func (a *AgentState) ReceiveMsg() (string, bool) {
	select {
	case msg := <-a.Mail:
		return msg, true
	default:
		return "", false
	}
}

//ClearMail throws away any message on the Agent's Mail channel
func (a *AgentState) ClearMail() {
	select {
	case <-a.Mail:
		return
	default:
		return
	}
}
