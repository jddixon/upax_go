package upax_go

// upax_go/s_msg_handlers.go

// Message handlers for messages betweeen Upax servers, that is, for
// intra-cluster communications.

import (
	"bytes"
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
	"errors"
	"fmt"
	// xc "github.com/jddixon/xlattice_go/crypto"
	// xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
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

/////////////////////////////////////////////////////////////////////
// AES-BASED MESSAGE PAIRS
// All of these functions have the same signature, so that they can
// be invoked through a dispatch table.
/////////////////////////////////////////////////////////////////////

// Dispatch table entry where a client message received is inappropriate
// the the state of the connection.  For example, ...
//
func badCCombo(h *ClientInHandler) {
	h.errOut = reg.RcvdInvalidMsgForState
}

// 0. INTRO AND ACK ================================================

func doCIntroMsg(h *ClientInHandler) {

	// XXX THIS IS JUST A COPY OF doCItsMe -- needs to be edited
	// considerably !

	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		clientMsg  *UpaxClientMsg
		clientID   []byte
		clientInfo *reg.MemberInfo
	)
	// expect clientMsgN to be 1
	err = checkCMsgN(h)
	if err == nil {
		clientMsg = h.msgIn

		// XXX THIS IS SIMPLY WRONG - NO SUCH FIELD
		clientID = clientMsg.GetID()

		// use the clientID to get their memberInfo
		for i := 0; i < len(h.us.Members); i++ {
			memberInfo := h.us.Members[i]
			if bytes.Equal(clientID, memberInfo.GetNodeID().Value()) {
				clientInfo = memberInfo
				break
			}
		}
		if h.clientInfo == nil {
			err = UnknownClient
		}
	}
	if err == nil {
		peerSK := h.clientInfo.GetSigPublicKey()
		salt := clientMsg.GetSalt()
		sig := clientMsg.GetSig()

		// use the public key to verify their digsig on the fields
		// present in canonical order: id, salt
		if sig == nil {
			err = NoDigSig
		} else {
			if clientID == nil && salt == nil {
				err = NoSigFields
			} else {
				d := sha1.New()
				if clientID != nil {
					d.Write(clientID)
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
		// the ClientInHandler.
		h.clientInfo = clientInfo

	}
	if err == nil {
		// Send reply to client -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
}

// 1. ITS_ME AND ACK ================================================

// Handle an ItsMe msg: we return an Ack or closes the connection.
// This should normally take the connection to an C_ID_VERIFIED state.
//
func doCItsMeMsg(h *ClientInHandler) {

	// XXX ALL peer TO BE REPLACED BY client

	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		clientMsg  *UpaxClientMsg
		clientID   []byte
		clientInfo *reg.MemberInfo
	)
	// expect clientMsgN to be 1
	err = checkCMsgN(h)
	if err == nil {
		clientMsg = h.msgIn
		clientID = clientMsg.GetID()

		// use the clientID to get their memberInfo
		for i := 0; i < len(h.us.Members); i++ {
			memberInfo := h.us.Members[i]
			if bytes.Equal(clientID, memberInfo.GetNodeID().Value()) {
				clientInfo = memberInfo
				break
			}
		}
		if h.clientInfo == nil {
			err = UnknownClient
		}
	}
	if err == nil {
		peerSK := h.clientInfo.GetSigPublicKey()
		salt := clientMsg.GetSalt()
		sig := clientMsg.GetSig()

		// use the public key to verify their digsig on the fields
		// present in canonical order: id, salt
		if sig == nil {
			err = NoDigSig
		} else {
			if clientID == nil && salt == nil {
				err = NoSigFields
			} else {
				d := sha1.New()
				if clientID != nil {
					d.Write(clientID)
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
		// the ClientInHandler.
		h.clientInfo = clientInfo

	}
	if err == nil {
		// Send reply to client -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
} // GEEP

// 2. KEEP-ALIVE AND ACK ============================================

// Handle a KeepAlive msg: we just return an Ack

func doCKeepAliveMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	err = checkCMsgN(h)

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED // the base state
	}
}

// 3. QUERY AND ACK ================================================

//
func doCQueryMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		queryMsg *UpaxClientMsg
	)
	err = checkCMsgN(h)
	if err == nil {
		queryMsg = h.msgIn
	}
	if err == nil {
	}
	_ = queryMsg

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
} // GEEP

// 4. GET AND DATA ==================================================

// Handle a Get msg.  If we have the data, we return it as a DataMsg
// (payload plus log entry); otherwise we will return NotFound, a non-fatal
// error message.

func doCGetMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		getMsg *UpaxClientMsg
		found  bool
	)
	err = checkCMsgN(h)
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

// 5. PUT AND ACK  ==================================================

//
func doCPutMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		putMsg *UpaxClientMsg
	)
	err = checkCMsgN(h)
	if err == nil {
		putMsg = h.msgIn
	}
	_ = putMsg

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client ----------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
}

// 6. BYE AND ACK ===================================================

func doCByeMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()

	// Examine incoming message -------------------------------------
	err = checkCMsgN(h)

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
}
