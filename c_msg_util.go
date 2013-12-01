package upax_go

// upax_go/c_msg_util.go

// Utility functions for messages betweeen Upax servers.

import (
	"errors"
	"fmt"
)

// Verify that the message number on the incoming message has been
// increased by one.
//
func checkCMsgN(h *ClientInHandler) (err error) {
	byeMsg := h.msgIn
	peerMsgN := byeMsg.GetMsgN()
	expectedMsgN := h.peerMsgN + 1
	if peerMsgN != expectedMsgN {
		msg := fmt.Sprintf("expected MsgN %d, got %d",
			expectedMsgN, peerMsgN)
		err = errors.New(msg)
	} else {
		h.peerMsgN++
	}
	return
}
func sendCAck(h *ClientInHandler) {
	h.myMsgN++
	op := UpaxClientMsg_Ack
	h.msgOut = &UpaxClientMsg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
}
func sendCNotFound(h *ClientInHandler) {
	h.myMsgN++
	op := UpaxClientMsg_NotFound
	h.msgOut = &UpaxClientMsg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
}
