package ${pkgName}

// ${pkgName}/${filePrefix}msg_util.go

// Utility functions for messages betweeen Upax servers.

import (
	"errors"
	"fmt"
)

// Verify that the message number on the incoming message has been
// increased by one.
//
func check${CapShortPrefix}MsgN(h *${TypePrefix}InHandler) (err error) {
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
func send${CapShortPrefix}Ack(h *${TypePrefix}InHandler) {
	h.myMsgN++
	op := Upax${TypePrefix}Msg_Ack
	h.msgOut = &Upax${TypePrefix}Msg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
}
func send${CapShortPrefix}NotFound(h *${TypePrefix}InHandler) {
	h.myMsgN++
	op := Upax${TypePrefix}Msg_NotFound
	h.msgOut = &Upax${TypePrefix}Msg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
} 
