package orgnetsim

import (
	"testing"
)

type testNetwork struct {
	relatedAgents []*Agent
	agentByID     map[string]*Agent
}

func (tn *testNetwork) GetRelatedAgents(a *Agent) []*Agent {
	return tn.relatedAgents
}

func (tn *testNetwork) GetAgentByID(id string) *Agent {
	return tn.agentByID[id]
}

func (tn *testNetwork) UpdateLinkStrength(id1 string, id2 string) error {
	return nil
}

func (tn *testNetwork) addAgent(a *Agent) {
	tn.relatedAgents = append(tn.relatedAgents, a)
	tn.agentByID[a.ID] = a
}

func newTestNetwork() *testNetwork {
	tn := testNetwork{}
	tn.relatedAgents = []*Agent{}
	tn.agentByID = map[string]*Agent{}
	ra1 := newAgent()
	ra1.ID = "id_1"
	ra1.Influence = 0.5
	ra1.Susceptability = 0.5
	ra1.Contrariness = 0.5
	ra1.Color = Blue
	ra1.Matched <- true
	tn.addAgent(ra1)
	ra2 := newAgent()
	ra2.ID = "id_2"
	ra2.Matched <- true
	tn.addAgent(ra2)
	ra3 := newAgent()
	ra3.ID = "id_3"
	tn.addAgent(ra3)
	return &tn
}

func newAgent() *Agent {
	a := Agent{}
	a.Mail = make(chan string, 1)
	a.Matched = make(chan bool, 1)
	return &a
}

func TestInteractSendsMsgToFirstAvailableRelatedAgent(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Interact(tn)
	matched := <-tn.GetAgentByID("id_3").Matched
	IsTrue(t, matched, "First free agent in related agents is not matched")
	msg, received := tn.GetAgentByID("id_3").RecieveMsg()
	IsTrue(t, received, "Message not sent to first free agent in related agents is not matched")
	AreEqual(t, "id_aut", msg, "Wrong message sent to first free agent")
}

func TestInteractReceivesMsgHigherSusceptabilityIgnoresMsg(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Susceptability = 1
	aut.Color = Red
	sent := aut.SendMsg("id_1")
	aut.Interact(tn)
	IsTrue(t, sent, "Msg not sent to Agent under test")
	AreEqual(t, Red, aut.Color, "Agent Color should not change if Agent has higher susceptability")
}

func TestInteractReceivesMsgLowerSusceptabilityChangesColor(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Susceptability = 0.4
	aut.Color = Red
	sent := aut.SendMsg("id_1")
	aut.Interact(tn)
	IsTrue(t, sent, "Msg not sent to Agent under test")
	AreEqual(t, Blue, aut.Color, "Agent Color should change to Blue if Agent has lower susceptability")
}

func TestInteractReceivesMsgLowerSusceptabilityHigherContrarinessRadomlyChangesColor(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Susceptability = 0.4
	aut.Contrariness = 0.6
	aut.Color = Red
	sent := aut.SendMsg("id_1")
	aut.Interact(tn)
	IsTrue(t, sent, "Msg not sent to Agent under test")
	NotEqual(t, Blue, aut.Color, "Agent Color should change to random Color if Agent has higher contrariness")
	NotEqual(t, Red, aut.Color, "Agent Color should change to random Color if Agent has higher contrariness")
}

func TestRecieveMsgTimesOutWhenNoMsg(t *testing.T) {
	a := newAgent()
	_, received := a.RecieveMsg()
	IsFalse(t, received, "Unexpected true returned from ReceiveMsg")
}

func TestRecieveMsgGetsMsg(t *testing.T) {
	a := newAgent()
	origMsg := "myMsg"
	a.Mail <- origMsg
	msg, received := a.RecieveMsg()
	IsTrue(t, received, "Unexpected false returned from ReceiveMsg")
	AreEqual(t, origMsg, msg, "Msgs not equal.")
}

func TestSendMsgSendsMsg(t *testing.T) {
	a := newAgent()
	origMsg := "myMsg"
	sent := a.SendMsg(origMsg)
	IsTrue(t, sent, "Msg not sent")
}

func TestSendMsgFailsSecondTime(t *testing.T) {
	a := newAgent()
	origMsg := "myMsg"
	a.SendMsg(origMsg)
	sent := a.SendMsg(origMsg)
	IsFalse(t, sent, "2nd msg sent but should have returned failed and returned false")
}

func TestClearMsg(t *testing.T) {
	a := newAgent()
	origMsg := "myMsg"
	a.Mail <- origMsg
	a.ClearMsg()
	select {
	case msg := <-a.Mail:
		t.Errorf("%s not cleared from Mail channel", msg)
		return
	default:
		return
	}
}

func TestClearMsgDoesntBlockWhenChannelEmpty(t *testing.T) {
	a := newAgent()
	a.ClearMsg()
}

func TestTryMatchReturnsTrueWhenMatched(t *testing.T) {
	a := newAgent()
	matched := a.TryMatch()
	IsTrue(t, matched, "Agent not matched but should have been")
}

func TestTryMatchReturnsFalseWhenBlocked(t *testing.T) {
	a := newAgent()
	a.TryMatch()
	matched := a.TryMatch()
	IsFalse(t, matched, "Agent matched but should have been blocked")
}

func TestClearMatch(t *testing.T) {
	a := newAgent()
	a.Matched <- true
	a.ClearMatch()
	select {
	case matched := <-a.Matched:
		t.Errorf("%t not cleared from Mail channel", matched)
		return
	default:
		return
	}
}

func TestClearInteractionsClearsMailAndMatchedChannels(t *testing.T) {
	a := newAgent()
	a.Matched <- true
	a.Mail <- "msg"
	a.ClearInteractions()
	select {
	case matched := <-a.Matched:
		t.Errorf("%t not cleared from Matched channel", matched)
		return
	case msg := <-a.Mail:
		t.Errorf("%s not cleared from Mail channel", msg)
		return
	default:
		return
	}
}

func TestClearMatchDoesntBlockWhenChannelEmpty(t *testing.T) {
	a := newAgent()
	a.ClearMatch()
}

func TestDefaultColorIsGrey(t *testing.T) {
	a := newAgent()
	AreEqual(t, a.Color, Grey, "Default Agent color is not Grey")
}

func TestSetColorChangesColor(t *testing.T) {
	a := newAgent()
	a.SetColor(Blue)
	AreEqual(t, a.Color, Blue, "Agent color is not set to Blue")
}

func TestSetColorIncrementsChangeCount(t *testing.T) {
	a := newAgent()
	AreEqual(t, a.ChangeCount, int32(0), "Default change count is not set to 0")
	a.SetColor(Blue)
	AreEqual(t, a.ChangeCount, int32(1), "Change count is not incremented to 1")
	a.SetColor(Red)
	AreEqual(t, a.ChangeCount, int32(2), "Change count is not incremented to 2")
}
