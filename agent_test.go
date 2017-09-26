package orgnetsim

import "testing"

func TestRecieveMsgTimesOutWhenNoMsg(t *testing.T) {
	a := Agent{}
	a.Mail = make(chan string, 1)
	_, received := a.RecieveMsg()
	if received {
		t.Error("Unexpected true returned from ReceiveMsg")
	}
	return
}

func TestRecieveMsgGetsMsg(t *testing.T) {
	a := Agent{}
	a.Mail = make(chan string, 1)
	origMsg := "myMsg"
	a.Mail <- origMsg
	msg, received := a.RecieveMsg()
	if !received {
		t.Error("Unexpected false returned from ReceiveMsg")
	}
	if msg != origMsg {
		t.Errorf("Msgs not equal. Expected %s got %s", origMsg, msg)
	}
	return
}

func TestSendMsgSendsMsg(t *testing.T) {
	a := Agent{}
	a.Mail = make(chan string, 1)
	origMsg := "myMsg"
	sent := a.SendMsg(origMsg)
	if !sent {
		t.Error("Msg not sent")
	}
	return
}

func TestSendMsgFailsSecondTime(t *testing.T) {
	a := Agent{}
	a.Mail = make(chan string, 1)
	origMsg := "myMsg"
	a.SendMsg(origMsg)
	sent := a.SendMsg(origMsg)
	if sent {
		t.Error("2nd msg sent but should have returned failed and returned false")
	}
	return
}

func TestClearMsg(t *testing.T) {
	a := Agent{}
	a.Mail = make(chan string, 1)
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
	a := Agent{}
	a.Mail = make(chan string, 1)
	a.ClearMsg()
	return
}

func TestTryMatchReturnsTrueWhenMatched(t *testing.T) {
	a := Agent{}
	a.Matched = make(chan bool, 1)
	matched := a.TryMatch()
	if !matched {
		t.Error("Agent not matched but should have been")
	}
	return
}

func TestTryMatchReturnsFalseWhenBlocked(t *testing.T) {
	a := Agent{}
	a.Matched = make(chan bool, 1)
	a.TryMatch()
	matched := a.TryMatch()
	if matched {
		t.Error("Agent matched but should have been blocked")
	}
	return
}

func TestClearMatch(t *testing.T) {
	a := Agent{}
	a.Matched = make(chan bool, 1)
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
	a := Agent{}
	a.Matched = make(chan bool, 1)
	a.ClearMatch()
	return
}

func TestDefaultColorIsGrey(t *testing.T) {
	a := Agent{}
	if a.Color != Grey {
		t.Error("Default Agent color is not Grey")
	}
}

func TestSetColorChangesColor(t *testing.T) {
	a := Agent{}
	a.SetColor(Blue)
	if a.Color != Blue {
		t.Error("Agent color is not set to Blue")
	}
}

func TestSetColorIncrementsChangeCount(t *testing.T) {
	a := Agent{}
	if a.ChangeCount != 0 {
		t.Error("Default change count is not set to 0")
	}
	a.SetColor(Blue)
	if a.ChangeCount != 1 {
		t.Error("Change count is not incremented to 1")
	}
	a.SetColor(Red)
	if a.ChangeCount != 2 {
		t.Error("Change count is not incremented to 2")
	}
}
