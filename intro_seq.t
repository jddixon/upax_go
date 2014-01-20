package ${pkgName}

// ${pkgName}/c_intro_seq.go

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
func do${CapShortPrefix}IntroMsg(h *${TypePrefix}InHandler) {

	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		peerMsg                         *Upax${TypePrefix}Msg
		name                              string
		token                             *Upax${TypePrefix}Msg_Token
		rawID, ckRaw, skRaw, salt, digSig []byte
		peerCK, peerSK                *rsa.PublicKey
		peerID                          *xi.NodeID
		peerInfo                        *reg.MemberInfo
	)
	// expect peerMsgN to be 1
	err = check${CapShortPrefix}MsgN(h)
	if err == nil {
		peerMsg = h.msgIn
		token = peerMsg.Get${TypePrefix}Info()
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
		peerID, err = xi.New(rawID)
		if err == nil {
			peerCK, err = xc.RSAPubKeyFromWire(ckRaw)
		}
		if err == nil {
			peerSK, err = xc.RSAPubKeyFromWire(ckRaw)
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
		err = rsa.VerifyPKCS1v15(peerSK, cr.SHA1, hash, digSig)
	}
	if err == nil {
		peerInfo, err = reg.NewMemberInfo(name, peerID,
			peerCK, peerSK, 0, nil)
	}
	// Take appropriate action --------------------------------------
	if err == nil {
		// The appropriate action is to hang a token for this client off
		// the ${TypePrefix}InHandler.
		h.peerInfo = peerInfo
	}
	if err == nil {
		// Send reply to client -------------------------------------
		send${CapShortPrefix}Ack(h)

		// Set exit state -------------------------------------------
		h.exitState = ${ConstPrefix}ID_VERIFIED
	}
}