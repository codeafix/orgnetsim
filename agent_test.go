package orgnetsim

import (
	"testing"
)

func NewAgent() *Agent {
	a := Agent{}
	a.Mail = make(chan string, 1)
	a.Matched = make(chan bool, 1)
	return &a
}

func TestRecieveMsgTimesOutWhenNoMsg(t *testing.T) {
	a := NewAgent()
	_, received := a.RecieveMsg()
	IsFalse(t, received, "Unexpected true returned from ReceiveMsg")
}

func TestRecieveMsgGetsMsg(t *testing.T) {
	a := NewAgent()
	origMsg := "myMsg"
	a.Mail <- origMsg
	msg, received := a.RecieveMsg()
	IsTrue(t, received, "Unexpected false returned from ReceiveMsg")
	AreEqual(t, origMsg, msg, "Msgs not equal.")
}

func TestSendMsgSendsMsg(t *testing.T) {
	a := NewAgent()
	origMsg := "myMsg"
	sent := a.SendMsg(origMsg)
	IsTrue(t, sent, "Msg not sent")
}

func TestSendMsgFailsSecondTime(t *testing.T) {
	a := NewAgent()
	origMsg := "myMsg"
	a.SendMsg(origMsg)
	sent := a.SendMsg(origMsg)
	IsFalse(t, sent, "2nd msg sent but should have returned failed and returned false")
}

func TestClearMsg(t *testing.T) {
	a := NewAgent()
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
	a := NewAgent()
	a.ClearMsg()
}

func TestTryMatchReturnsTrueWhenMatched(t *testing.T) {
	a := NewAgent()
	matched := a.TryMatch()
	IsTrue(t, matched, "Agent not matched but should have been")
}

func TestTryMatchReturnsFalseWhenBlocked(t *testing.T) {
	a := NewAgent()
	a.TryMatch()
	matched := a.TryMatch()
	IsFalse(t, matched, "Agent matched but should have been blocked")
}

func TestClearMatch(t *testing.T) {
	a := NewAgent()
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

func TestClearMatchDoesntBlockWhenChannelEmpty(t *testing.T) {
	a := NewAgent()
	a.ClearMatch()
}

func TestDefaultColorIsGrey(t *testing.T) {
	a := NewAgent()
	AreEqual(t, a.Color, Grey, "Default Agent color is not Grey")
}

func TestSetColorChangesColor(t *testing.T) {
	a := NewAgent()
	a.SetColor(Blue)
	AreEqual(t, a.Color, Blue, "Agent color is not set to Blue")
}

func TestSetColorIncrementsChangeCount(t *testing.T) {
	a := NewAgent()
	AreEqual(t, a.ChangeCount, int32(0), "Default change count is not set to 0")
	a.SetColor(Blue)
	AreEqual(t, a.ChangeCount, int32(1), "Change count is not incremented to 1")
	a.SetColor(Red)
	AreEqual(t, a.ChangeCount, int32(2), "Change count is not incremented to 2")
}
