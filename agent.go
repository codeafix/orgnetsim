package orgnetsim

import "time"

//An Agent is a node in the network
type Agent struct {
	ID             string      `json:"id"`
	Color          Color       `json:"color"`
	Susceptability float64     `json:"susceptability"`
	Influence      float64     `json:"influence"`
	Contrariness   float64     `json:"contrariness"`
	Matched        chan bool   `json:"-"`
	Mail           chan string `json:"-"`
	ChangeCount    int         `json:"change"`
}

//AgentInteracter is an interface that allows interaction with an Agent
type AgentInteracter interface {
	SendMail(n RelationshipMgr)
	ReadMail(n RelationshipMgr)
	ClearInteractions()
}

/*SendMail iterates over a randomly ordered slice of related agents trying to find a match. It sends a mail to the
first successful match it finds.
*/
func (a *Agent) SendMail(n RelationshipMgr) {
	for _, ra := range n.GetRelatedAgents(a) {
		if ra.TryMatch() {
			ra.SendMsg(a.ID)
			break
		} else {
			continue
		}
	}
}

/*ReadMail checks for any messages it recieved in its own Mail queue. If it receives
one then it decides whether to update its color.
*/
func (a *Agent) ReadMail(n RelationshipMgr) {
	msg, received := a.RecieveMsg()
	if received {
		ra := n.GetAgentByID(msg)
		n.IncrementLinkStrength(a.ID, ra.ID)
		if ra.Influence > a.Susceptability {
			if a.Contrariness > ra.Influence {
				altColor := RandomlySelectAlternateColor(ra.Color)
				a.SetColor(altColor)
			} else {
				a.SetColor(ra.Color)
			}
		}
	}
}

// ClearMail clears the matched and mail channels
func (a *Agent) ClearMail() {
	a.ClearMsg()
	a.ClearMatch()
}

//SetColor changes the color of the current Agent and counts the number of times the Agent changes color
func (a *Agent) SetColor(color Color) {
	if a.Color != color {
		a.ChangeCount = a.ChangeCount + 1
		a.Color = color
	}
}

//TryMatch tries to add an entry into an Agent's Matched channel, if it succeeds, that Agent will be blocked
//for matching to any other Agent and this function returns true (the Agent is Matched). If it returns
//false the Agent is already matched.
func (a *Agent) TryMatch() bool {
	select {
	case a.Matched <- true:
		return true
	default:
		return false
	}
}

//SendMsg adds a message to the Agent's Mail channel
func (a *Agent) SendMsg(msg string) bool {
	select {
	case a.Mail <- msg:
		return true
	default:
		return false
	}
}

//RecieveMsg picks a message up from the Agent's Mail channel
func (a *Agent) RecieveMsg() (string, bool) {
	select {
	case msg := <-a.Mail:
		return msg, true
	case <-time.After(100 * time.Millisecond):
		return "", false
	}
}

//ClearMsg throws away any message on the Agent's Mail channel
func (a *Agent) ClearMsg() {
	select {
	case <-a.Mail:
		return
	default:
		return
	}
}

//ClearMatch clears the flag on the Agent's Matched channel
func (a *Agent) ClearMatch() {
	select {
	case <-a.Matched:
		return
	default:
		return
	}
}
