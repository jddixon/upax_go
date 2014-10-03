package upax_go

// upax_go/c_msg_handlers.go

// Message handlers for messages betweeen Upax servers, that is, for
// intra-cluster communications.

import (
	"bytes"
	cr "crypto"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/hex"
	reg "github.com/jddixon/xlReg_go"
)

/////////////////////////////////////////////////////////////////////
// AES-BASED MESSAGE PAIRS
// All of these functions have the same signature, so that they can
// be invoked through a dispatch table.
/////////////////////////////////////////////////////////////////////

// 1. ITS_ME AND ACK ================================================

// Handle an ItsMe msg: we return an Ack or closes the connection.
// This should normally take the connection to an C_ID_VERIFIED state.
//
func doCItsMeMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		peerMsg  *UpaxClientMsg
		peerID   []byte
		peerInfo *reg.MemberInfo
	)
	// expect peerMsgN to be 1
	err = checkCMsgN(h)
	if err == nil {
		peerMsg = h.msgIn
		peerID = peerMsg.GetID()

		// use the peerID to get their memberInfo
		for i := 0; i < len(h.us.Members); i++ {
			memberInfo := h.us.Members[i]
			if bytes.Equal(peerID, memberInfo.GetNodeID().Value()) {
				peerInfo = memberInfo
				break
			}
		}
		if h.peerInfo == nil {
			err = UnknownPeer
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
		// The appropriate action is to hang a token for this peer off
		// the ClientInHandler.
		h.peerInfo = peerInfo

	}
	if err == nil {
		// Send reply to peer -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
}

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
		// Send reply to peer -------------------------------------
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
		// Send reply to peer -------------------------------------
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
		found     bool
		getMsg    *UpaxClientMsg
		data, key []byte
		logEntry  *LogEntry // not the message
	)
	err = checkCMsgN(h)
	if err == nil {
		getMsg = h.msgIn
		key = getMsg.GetHash()
		if key == nil {
			err = NilHash
		}
	}

	// Take appropriate action --------------------------------------
	if err == nil {
		// determine whether the data requested is present; if it is
		// we will send a DataMsg, with logEntry and payload fields
		var whatever interface{}
		whatever, err = h.us.entries.Find(key)
		logEntry = whatever.(*LogEntry)
		if err != nil {
			found = logEntry != nil
		}
	}
	if err == nil {
		if found {
			// fetch payload
			data, err = h.us.uDir.GetData(key)
			if err == nil {
				// we will send log entry and payload
				logEntryMsg := &UpaxClientMsg_LogEntry{
					Timestamp:  &logEntry.timestamp,
					ContentKey: logEntry.key,
					Owner:      logEntry.nodeID,
					Src:        &logEntry.src,
					Path:       &logEntry.path,
				}
				h.myMsgN++
				op := UpaxClientMsg_Data
				h.msgOut = &UpaxClientMsg{
					Op:      &op,
					MsgN:    &h.myMsgN,
					Entry:   logEntryMsg,
					Payload: data,
				}
			}

		} else {
			sendCNotFound(h)
		}
		if err == nil {
			// Set exit state -----------------------------------------------
			h.exitState = C_ID_VERIFIED
		}
	}
}

// 5. I_HAVE AND ACK ================================================

//
func doCIHaveMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		iHaveMsg *UpaxClientMsg
	)
	err = checkCMsgN(h)
	if err == nil {
		iHaveMsg = h.msgIn
	}
	if err == nil {

		// XXX STUB XXX - check the list, filtering out those we
		// already have; the list is then forwarded to the other
		// side (the outHandler), which will get anything on the list.
	}
	_ = iHaveMsg

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to peer -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_ID_VERIFIED
	}
}

// 6. PUT AND ACK  ==================================================

//
func doCPutMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()
	// Examine incoming message -------------------------------------
	var (
		data, key []byte
		entryMsg  *UpaxClientMsg_LogEntry
		logEntry  *LogEntry
		putMsg    *UpaxClientMsg
	)
	err = checkCMsgN(h)
	if err == nil {
		putMsg = h.msgIn
		entryMsg = putMsg.GetEntry()
		if entryMsg == nil {
			err = NilLogEntry
		}
	}
	if err == nil {
		data = putMsg.GetPayload()
		if data == nil {
			err = NilPayload
		}
	}
	if err == nil {
		t := entryMsg.GetTimestamp()
		key = entryMsg.GetContentKey()
		nodeID := entryMsg.GetOwner()
		src := entryMsg.GetSrc()
		path := entryMsg.GetPath()

		// XXX CHECK FOR MISSING FIELDS

		logEntry, err = NewLogEntry(t, key, nodeID, src, path)
	}
	// Take appropriate action --------------------------------------
	if err == nil {
		var (
			found    bool
			whatever interface{}
		)
		whatever, err = h.us.entries.Find(key)
		if err == nil {
			if whatever != nil {
				found = logEntry.Equal(whatever)
			}
			if !found {
				// write data to U/x/x, the appropriate hash directory
				var hash []byte
				_, hash, err = h.us.uDir.PutData(data, key)
				if err == nil && !bytes.Equal(hash, key) {
					err = BadHash
				}
			}
		}
		if err == nil && !found {
			// write to U/L, the log file (UNSYNCHRONIZED)
			_, err = h.us.ftLogFile.WriteString(logEntry.String() + "\n")
			if err == nil {
				// written to U, entry appended to U/L, so put it into the map
				err = h.us.entries.Insert(key, logEntry)
			}
		}
		if err == nil {
			// Send reply to peer ----------------------------------
			sendCAck(h)

			// Set exit state ---------------------------------------
			h.exitState = C_ID_VERIFIED
		}
	}
}

// 7. BYE AND ACK ===================================================

func doCByeMsg(h *ClientInHandler) {
	var err error
	defer func() {
		h.errOut = err
	}()

	// Examine incoming message -------------------------------------
	err = checkCMsgN(h)

	// Take appropriate action --------------------------------------
	if err == nil {
		// Send reply to peer -------------------------------------
		sendCAck(h)

		// Set exit state -------------------------------------------
		h.exitState = C_BYE_RCVD
	}
}
