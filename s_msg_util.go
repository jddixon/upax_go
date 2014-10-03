package upax_go

// upax_go/s_msg_util.go

// Utility functions for messages betweeen Upax servers.

import (
	"errors"
	"fmt"
	reg "github.com/jddixon/xlReg_go"
)

// Verify that the message number on the incoming message has been
// increased by one.
//
func checkSMsgN(h *ClusterInHandler) (err error) {
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
func sendSAck(h *ClusterInHandler) {
	h.myMsgN++
	op := UpaxClusterMsg_Ack
	h.msgOut = &UpaxClusterMsg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
}
func sendSNotFound(h *ClusterInHandler) {
	h.myMsgN++
	op := UpaxClusterMsg_NotFound
	h.msgOut = &UpaxClusterMsg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
}

// Dispatch table entry where a message received is inappropriate
// the the state of the connection.  For example, ...
func badSCombo(h *ClusterInHandler) {
	h.errOut = reg.RcvdInvalidMsgForState
}
