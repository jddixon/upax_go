package upax_go

// upax_go/s_msg_handlers.go

// Message handlers for messages betweeen Upax servers, that is, for
// intra-cluster communications.

import (
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
	"errors"
	"fmt"
	// xc "github.com/jddixon/xlattice_go/crypto"
	// xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
	xu "github.com/jddixon/xlattice_go/util"
)

// Verify that the message number on the incoming message has been
// increased by one.
//
func checkMsgN(h *ClusterInHandler) (err error) {
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
func sendAck(h *ClusterInHandler) {
	h.myMsgN++
	op := UpaxClusterMsg_Ack
	h.msgOut = &UpaxClusterMsg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
}
func sendNotFound(h *ClusterInHandler) {
	h.myMsgN++
	op := UpaxClusterMsg_NotFound
	h.msgOut = &UpaxClusterMsg{
		Op:       &op,
		MsgN:     &h.myMsgN,
		YourMsgN: &h.peerMsgN,
	}
}

/////////////////////////////////////////////////////////////////////
// AES-BASED MESSAGE PAIRS
// All of these functions have the same signature, so that they can
// be invoked through a dispatch table.
/////////////////////////////////////////////////////////////////////

// Dispatch table entry where a client message received is inappropriate
// the the state of the connection.  For example, ...
//
func badCombo(h *ClusterInHandler) {
	h.errOut = reg.RcvdInvalidMsgForState
}

// 0. ITS_ME AND ACK ================================================

// Handle an ItsMe msg: we return an Ack or closes the connection.
// This should normally take the connection to an ID_VERIFIED state.
//
func doItsMeMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		peerMsg  *UpaxClusterMsg
		peerID   []byte
		peerInfo *reg.MemberInfo
	)
	// expect peerMsgN to be 1
	err = checkMsgN(h)
	if err == nil {
		peerMsg = h.msgIn
		peerID = peerMsg.GetID()

		// use the peerID to get their memberInfo
		for i := 0; i < len(h.us.Members); i++ {
			memberInfo := h.us.Members[i]
			if xu.SameBytes(peerID, memberInfo.GetNodeID().Value()) {
				peerInfo = memberInfo
				break
			}
		}
		if h.peerInfo == nil {
			err = NotClusterMember
		}
	}
	if err == nil {
		peerSK := h.peerInfo.GetSigPublicKey()
		salt := peerMsg.GetSalt()
		sig := peerMsg.GetSig()

		// use the public key to verify their digsig on the fields
		// present in canonical order: id, salt
		if sig == nil {
			err = NoDigSig
		} else {
			if peerID == nil && salt == nil {
				err = NoSigFields
			} else {
				d := sha1.New()
				if peerID != nil {
					d.Write(peerID)
				}
				if salt != nil {
					d.Write(salt)
				}
				hash := d.Sum(nil)
				err = rsa.VerifyPKCS1v15(peerSK, cr.SHA1, hash, sig)
			}
		}
	}
	// Take appropriate action --------------------------------------
	if err == nil {
		// The appropriate action is to hang a token for this client off
		// the ClusterInHandler.
		h.peerInfo = peerInfo

	}
	if err == nil {
		// Send reply to client -------------------------------------
		sendAck(h)

		// Set exit state -------------------------------------------
		h.exitState = ID_VERIFIED
	}
}

// 1. KEEP-ALIVE AND ACK ============================================

// Handle a KeepAlive msg: we just return an Ack

func doKeepAliveMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	err = checkMsgN(h)

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client -------------------------------------
		sendAck(h)

		// Set exit state -------------------------------------------
		h.exitState = ID_VERIFIED // the base state
	}
}

// 2. GET AND DATA ==================================================

// Handle a Get msg.  If we have the data, we return it as a DataMsg
// (payload plus log entry); otherwise we will return NotFound, a non-fatal
// error message.

func doGetMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		getMsg *UpaxClusterMsg
		found  bool
	)
	err = checkMsgN(h)
	if err == nil {
		getMsg = h.msgIn
	}
	_, _ = found, getMsg // DEBUG

	// Take appropriate action --------------------------------------
	if err == nil {
		// determine whether the data requested is present; if it is
		// we will send a DataMsg, with logEntry and payload fields

		// if the data is not present, send NotFound

	}
	if err == nil {
		// Set exit state -----------------------------------------------
		// h.exitState = CREATE_REQUEST_RCVD
	}
}

// 3. I_HAVE AND ACK ================================================

//
func doIHaveMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		iHaveMsg *UpaxClusterMsg
	)
	err = checkMsgN(h)
	if err == nil {
		iHaveMsg = h.msgIn
	}
	if err == nil {
	}
	_ = iHaveMsg

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client -------------------------------------
		sendAck(h)

		// Set exit state -------------------------------------------
		h.exitState = ID_VERIFIED
	}
}

// 4. PUT AND ACK  ==================================================

//
func doPutMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		putMsg *UpaxClusterMsg
	)
	err = checkMsgN(h)
	if err == nil {
		putMsg = h.msgIn
	}
	_ = putMsg

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client ----------------------------------
		sendAck(h)

		// Set exit state -------------------------------------------
		h.exitState = ID_VERIFIED
	}
}

// 5. BYE AND ACK ===================================================

func doByeMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()

	// Examine incoming message -------------------------------------
	err = checkMsgN(h)

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client -------------------------------------
		sendAck(h)

		// Set exit state -------------------------------------------
		h.exitState = ID_VERIFIED
	}
}
