package orgnetsim

import (
	"testing"
)

type testNetwork struct {
	relatedAgents []*Agent
	agentByID     map[string]*Agent
	LinkStrength  int
}

func (tn *testNetwork) GetRelatedAgents(a *Agent) []*Agent {
	return tn.relatedAgents
}

func (tn *testNetwork) GetAgentByID(id string) *Agent {
	return tn.agentByID[id]
}

func (tn *testNetwork) IncrementLinkStrength(id1 string, id2 string) error {
	tn.LinkStrength = tn.LinkStrength + 1
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
	ra1.Mail <- "id_5"
	tn.addAgent(ra1)
	ra2 := newAgent()
	ra2.ID = "id_2"
	ra2.Mail <- "id_6"
	tn.addAgent(ra2)
	ra3 := newAgent()
	ra3.ID = "id_3"
	tn.addAgent(ra3)
	return &tn
}

func newAgent() *Agent {
	a := Agent{}
	a.Mail = make(chan string, 1)
	return &a
}

func TestSendMailSendsMsgToFirstAvailableRelatedAgent(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	count := aut.SendMail(tn)
	AreEqual(t, 1, count, "Message not sent to first free Agent")
	msg, received := tn.GetAgentByID("id_3").RecieveMsg()
	IsTrue(t, received, "Message not sent to first free agent in related agents")
	AreEqual(t, "id_aut", msg, "Wrong message sent to first free agent")
}

func TestSendMailDoesNotSendIfNoAvailableRelatedAgent(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	tn.GetAgentByID("id_3").Mail <- "block"
	count := aut.SendMail(tn)
	AreEqual(t, 0, count, "Message sent but should not have been since there are no Agents free")
}

func TestReadMailReceivesMsgIncrementsLinkStrength(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Susceptability = 1
	aut.Color = Red
	sent := aut.SendMsg("id_1")
	AreEqual(t, 0, tn.LinkStrength, "LinkStrength not initialised to 0")
	aut.ReadMail(tn)
	IsTrue(t, sent, "Msg not sent to Agent under test")
	AreEqual(t, 1, tn.LinkStrength, "LinkStrength not incremented as part of reading a Mail")
	sent = aut.SendMsg("id_1")
	IsTrue(t, sent, "Msg not sent to Agent under test")
	aut.ReadMail(tn)
	AreEqual(t, 2, tn.LinkStrength, "LinkStrength not incremented as part of reading a Mail")
}

func TestReadMailReceivesMsgHigherSusceptabilityIgnoresMsg(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Susceptability = 1
	aut.Color = Red
	sent := aut.SendMsg("id_1")
	aut.ReadMail(tn)
	IsTrue(t, sent, "Msg not sent to Agent under test")
	AreEqual(t, Red, aut.Color, "Agent Color should not change if Agent has higher susceptability")
}

func TestReadMailReceivesMsgLowerSusceptabilityChangesColor(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Susceptability = 0.4
	aut.Color = Red
	sent := aut.SendMsg("id_1")
	aut.ReadMail(tn)
	IsTrue(t, sent, "Msg not sent to Agent under test")
	AreEqual(t, Blue, aut.Color, "Agent Color should change to Blue if Agent has lower susceptability")
}

func TestReadMailReceivesMsgLowerSusceptabilityHigherContrarinessRadomlyChangesColor(t *testing.T) {
	tn := newTestNetwork()
	aut := newAgent()
	aut.ID = "id_aut"
	aut.Susceptability = 0.4
	aut.Contrariness = 0.6
	aut.Color = Red
	sent := aut.SendMsg("id_1")
	aut.ReadMail(tn)
	IsTrue(t, sent, "Msg not sent to Agent under test")
	NotEqual(t, Blue, aut.Color, "Agent Color should change to random Color if Agent has higher contrariness")
	NotEqual(t, Red, aut.Color, "Agent Color should change to random Color if Agent has higher contrariness")
}

func TestRecieveMsgReturnsFalseWhenNoMsg(t *testing.T) {
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

func TestClearMail(t *testing.T) {
	a := newAgent()
	origMsg := "myMsg"
	a.Mail <- origMsg
	a.ClearMail()
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
	a.ClearMail()
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
	AreEqual(t, a.ChangeCount, 0, "Default change count is not set to 0")
	a.SetColor(Blue)
	AreEqual(t, a.ChangeCount, 1, "Change count is not incremented to 1")
	a.SetColor(Red)
	AreEqual(t, a.ChangeCount, 2, "Change count is not incremented to 2")
}
