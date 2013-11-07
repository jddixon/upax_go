package upax_go

// upax_go/s_msg_handlers.go

// Message handlers for messages betweeen Upax servers, that is, for
// intra-cluster communications.

import (
	// "crypto/rsa"
	"errors"
	"fmt"
	// xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
	//xu "github.com/jddixon/xlattice_go/util"
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
		nodeID *xi.NodeID
	)
	// XXX We should accept EITHER clientName + token OR clientID
	// This implementation only accepts a token.

	peerMsg := h.msgIn
	peerMsgN := peerMsg.GetMsgN()
	peerID := peerMsg.GetID()
	salt := peerMsg.GetSalt()
	sig := peerMsg.GetSig()

	// expect peerMsgN to be 1
	err = checkMsgN(h)

	if err == nil {
		// use the peerID to get their public key

		// if not recognized,
		// err = NotClusterMember
	}
	if err == nil {
		// use the public key to verify their digsig on the fields
		// presesnt in canonical order: seqn, id, salt

		// if the digSig verification fails,
		// err = BadDigSig

	}
	// Take appropriate action --------------------------------------
	if err == nil {
		// The appropriate action is to hang a token for this client off
		// the ClusterInHandler.

	}
	if err == nil {
		// Prepare reply to client --------------------------------------
		op := UpaxClusterMsg_Ack
		h.msgOut = &UpaxClusterMsg{
			Op:   &op,
			MsgN: &h.myMsgN,
		}
		h.myMsgN++

		// Set exit state -----------------------------------------------
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

	if err == nil {

	}
	if err == nil {

	}
	// Take appropriate action --------------------------------------
	if err == nil {

	}
	if err == nil {
		// Prepare reply to client --------------------------------------
		op := UpaxClusterMsg_Ack
		h.msgOut = &UpaxClusterMsg{
			Op:   &op,
			MsgN: &h.myMsgN,
		}
		h.myMsgN++

		// Set exit state -----------------------------------------------
		h.exitState = ID_VERIFIED
	}
}

// 2. GET AND DATA ==================================================

// Handle a Get msg.  If we have the data, we return it as a DataMsg
// (payload plus log entry); otherwise we will return a non-fatal
// error message.

func doGetMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var ()
	getMsg := h.msgIn
	_ = getMsg

	// Take appropriate action --------------------------------------

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
	var ()
	iHaveMsg := h.msgIn
	_ = iHaveMsg

	// Take appropriate action --------------------------------------

	if err == nil {
	}
	if err == nil {
		// Prepare reply to client ----------------------------------
		// Set exit state -------------------------------------------
		// h.exitState = JOIN_RCVD
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
	var ()
	putMsg := h.msgIn
	_ = putMsg

	// Take appropriate action --------------------------------------
	if err == nil {
		// Prepare reply to client --------------------------------------

		// Set exit state -----------------------------------------------
		//h.exitState = JOIN_RCVD // the JOIN is intentional !
	}
}

// 5. BYE AND ACK ===================================================

func doByeMsg(h *ClusterInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()

	// Examine incoming message -------------------------------------
	byeMsg := h.msgIn

	// Take appropriate action --------------------------------------
	err = checkMsgN(h)

	if err == nil {
		// Prepare reply to client --------------------------------------
		op := UpaxClusterMsg_Ack
		h.msgOut = &UpaxClusterMsg{
			Op:   &op,
			MsgN: &h.myMsgN,
		}
		h.myMsgN++

		// Set exit state -----------------------------------------------
		h.exitState = BYE_RCVD
	}
}
