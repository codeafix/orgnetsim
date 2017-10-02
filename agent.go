package orgnetsim

//An Agent is a node in the network
type Agent struct {
	ID             string             `json:"id"`
	Color          Color              `json:"color"`
	Susceptability float64            `json:"susceptability"`
	Influence      float64            `json:"influence"`
	Contrariness   float64            `json:"contrariness"`
	Mail           chan string        `json:"-"`
	ChangeCount    int                `json:"change"`
	Memory         map[Color]struct{} `json:"-"`
}

//Communicator is an interface that allows interaction with an Agent
type Communicator interface {
	SendMail(n RelationshipMgr) int
	ReadMail(n RelationshipMgr) Color
	ClearMail()
	GetColor() Color
}

/*SendMail iterates over a randomly ordered slice of related agents trying to find a match. It sends a mail to the
first successful match it finds.
*/
func (a *Agent) SendMail(n RelationshipMgr) int {
	for _, ra := range n.GetRelatedAgents(a) {
		if ra.SendMsg(a.ID) {
			return 1
		}
	}
	return 0
}

/*ReadMail checks for any messages it received in its own Mail queue. If it receives
one then it decides whether to update its color.
*/
func (a *Agent) ReadMail(n RelationshipMgr) Color {
	msg, received := a.ReceiveMsg()
	if received {
		ra := n.GetAgentByID(msg)
		n.IncrementLinkStrength(a.ID, ra.ID)
		if ra.Influence > a.Susceptability {
			if a.Contrariness > ra.Influence {
				altColor := RandomlySelectAlternateColor(a.Color)
				a.SetColor(altColor)
			} else {
				a.SetColor(ra.Color)
			}
		}
	}
	return a.Color
}

//SetColor changes the color of the current Agent and counts the number of times the Agent changes color
//It also adds each color to a memory so that once it changes it's mind it doesn't change back
func (a *Agent) SetColor(color Color) {
	_, rem := a.Memory[color]
	if !rem && a.Color != color {
		a.ChangeCount++
		a.Color = color
		a.Memory[color] = struct{}{}
	}
}

//GetColor returns the Color of this Agent
func (a *Agent) GetColor() Color {
	return a.Color
}

//SendMsg tries to add an entry into an Agent's Mail channel, if it succeeds, that Agent will be blocked
//for any other Agent trying to send a Mail and this function returns true (the Agent is now Matched).
//If it returns false the Agent is already matched by another Agent.
func (a *Agent) SendMsg(msg string) bool {
	select {
	case a.Mail <- msg:
		return true
	default:
		return false
	}
}

//ReceiveMsg picks a message up from the Agent's Mail channel
func (a *Agent) ReceiveMsg() (string, bool) {
	select {
	case msg := <-a.Mail:
		return msg, true
	default:
		return "", false
	}
}

//ClearMail throws away any message on the Agent's Mail channel
func (a *Agent) ClearMail() {
	select {
	case <-a.Mail:
		return
	default:
		return
	}
}
