package upax_go

// upax_go/c_intro_seq.go

import (
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
)

/////////////////////////////////////////////////////////////////////
// AES-BASED MESSAGE PAIRS
// All of these functions have the same signature, so that they can
// be invoked through a dispatch table.
/////////////////////////////////////////////////////////////////////

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
