package upax_go

// upax_go/s_msg_handlers.go

// Message handlers for messages betweeen Upax servers, that is, for
// intra-cluster communications.

import (
	"bytes"
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
	xu "github.com/jddixon/xlattice_go/u"
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

// IntroMsg consists of MsgN and Token; the token should contain
// name, ID, commsKey, sigKey, salt, and a digital signature over
// fields present in the token, excluding the token itself.  On
// success the server replies with an IntroOK.
//
func doCIntroMsg(h *ClientInHandler) {

	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		clientMsg                         *UpaxClientMsg
		name                              string
		token                             *UpaxClientMsg_Token
		rawID, ckRaw, skRaw, salt, digSig []byte
		clientCK, clientSK                *rsa.PublicKey
		clientID                          *xi.NodeID
		clientInfo                        *reg.MemberInfo
	)
	// expect clientMsgN to be 1
	err = checkCMsgN(h)
	if err == nil {
		clientMsg = h.msgIn
		token = clientMsg.GetClientInfo()
		if token == nil {
			err = NilToken
		}
	}
	if err == nil {
		name = token.GetName()
		rawID = token.GetID()
		ckRaw = token.GetCommsKey()
		skRaw = token.GetSigKey()
		salt = token.GetSalt()
		digSig = token.GetDigSig()
		if name == "" || rawID == nil || ckRaw == nil || skRaw == nil ||
			salt == nil || digSig == nil {
			err = MissingTokenField
		}
	}
	if err == nil {
		clientID, err = xi.New(rawID)
		if err == nil {
			clientCK, err = xc.RSAPubKeyFromWire(ckRaw)
		}
		if err == nil {
			clientSK, err = xc.RSAPubKeyFromWire(ckRaw)
		}
	}
	if err == nil {
		// Use the public key to verify their digsig on the fields
		// present in canonical order

		d := sha1.New()
		d.Write([]byte(name))
		d.Write(rawID)
		d.Write(ckRaw)
		d.Write(skRaw)
		d.Write(salt)
		hash := d.Sum(nil)
		err = rsa.VerifyPKCS1v15(clientSK, cr.SHA1, hash, digSig)
	}
	if err == nil {
		clientInfo, err = reg.NewMemberInfo(name, clientID,
			clientCK, clientSK, 0, nil)
	}
	// Take appropriate action --------------------------------------
	if err == nil {
		// The appropriate action is to hang a token for this client off
		// the ClientInHandler.
		h.peerInfo = clientInfo
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

	// XXX ALL client TO BE REPLACED BY client

	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		clientMsg  *UpaxClientMsg
		rawID      []byte
		clientInfo *reg.MemberInfo
	)
	// expect clientMsgN to be 1
	err = checkCMsgN(h)
	if err == nil {
		clientMsg = h.msgIn
		rawID = clientMsg.GetID()

		// use the rawID to get their memberInfo
		for i := 0; i < len(h.us.Members); i++ {
			memberInfo := h.us.Members[i]
			if bytes.Equal(rawID, memberInfo.GetNodeID().Value()) {
				clientInfo = memberInfo
				break
			}
		}
		if h.peerInfo == nil {
			err = UnknownClient
		}
	}
	if err == nil {
		clientSK := h.peerInfo.GetSigPublicKey()
		salt := clientMsg.GetSalt()
		sig := clientMsg.GetSig()

		// use the public key to verify their digsig on the fields
		// present in canonical order: id, salt
		if sig == nil {
			err = NoDigSig
		} else {
			if rawID == nil && salt == nil {
				err = NoSigFields
			} else {
				d := sha1.New()
				if rawID != nil {
					d.Write(rawID)
				}
				if salt != nil {
					d.Write(salt)
				}
				hash := d.Sum(nil)
				err = rsa.VerifyPKCS1v15(clientSK, cr.SHA1, hash, sig)
			}
		}
	}
	// Take appropriate action --------------------------------------
	if err == nil {
		// The appropriate action is to hang a token for this client off
		// the ClientInHandler.
		h.peerInfo = clientInfo

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

// A faster implementation would check an in-memory Bloom filter,
// trusting a negative reply but verifying any positive with a disk hit.
//
func doCQueryMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		found    bool
		hash     []byte
		queryMsg *UpaxClientMsg
	)
	err = checkCMsgN(h)
	if err == nil {
		queryMsg = h.msgIn
		hash = queryMsg.GetHash()
		if hash == nil {
			err = NilHash
		}
	}
	if err == nil {
		strHash := hex.EncodeToString(hash)
		found, err = h.us.uDir.Exists(strHash)
	}

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to client -------------------------------------
		if found {
			sendCAck(h)
		} else {
			sendCNotFound(h)
		}

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
}

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
		getMsg    *UpaxClientMsg
		hash      []byte
		strHash   string
		usingSHA1 bool
	)
	err = checkCMsgN(h)
	if err == nil {
		getMsg = h.msgIn
		hash = getMsg.GetHash()
		if hash == nil {
			err = NilHash
		} else {
			strHash = hex.EncodeToString(hash)

			// BEWARE: U uses hex lengths, double byte lengths
			switch len(strHash) {
			case xu.SHA1_LEN:
				usingSHA1 = true
			case xu.SHA3_LEN:
				usingSHA1 = false
			default:
				err = BadHashLength
			}
		}
	}
	// Take appropriate action --------------------------------------
	if err == nil {
		var data []byte

		if usingSHA1 {
			data, err = h.us.uDir.GetData1(strHash)
		} else {
			data, err = h.us.uDir.GetData3(strHash)
		}
		if err == xu.FileNotFound {
			err = nil
			data = nil
		}
		if err == nil {
			if data == nil {
				sendCNotFound(h)
			} else {
				// we will send a DataMsg, with logEntry and payload fields

				h.myMsgN++
				op := UpaxClientMsg_Ack
				h.msgOut = &UpaxClientMsg{
					Op:   &op,
					MsgN: &h.myMsgN,

					// XXX LOG ENTRY MISSING!	<--------------------

					Payload: data,
				}

			}
		}
		// Set exit state -------------------------------------------
		if err == nil {
			h.exitState = C_ID_VERIFIED
		}
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
		h.exitState = C_BYE_RCVD
	}
}
