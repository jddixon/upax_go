package upax_go

// upax_go/client_msg.go

import (
	"crypto"
	"crypto/aes"
	//"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/binary"

	xr "github.com/jddixon/rnglib_go"
	xc "github.com/jddixon/xlCrypto_go"
	xa "github.com/jddixon/xlProtocol_go/aes_cnx"
	xt "github.com/jddixon/xlTransport_go"
	xu "github.com/jddixon/xlUtil_go"
)

// Read the next message over the connection
func (upc *UpaxClient) readMsg() (m *UpaxClientMsg, err error) {
	inBuf, err := upc.ReadData()
	if err == nil && inBuf != nil {
		m, err = clientDecryptUnpadDecode(inBuf, upc.decrypter)
	}
	return
}

// Write a message out over the connection
func (upc *UpaxClient) writeMsg(m *UpaxClientMsg) (err error) {
	var data []byte
	// serialize, marshal the message
	data, err = clientEncodePadEncrypt(m, upc.encrypter)
	if err == nil {
		err = upc.WriteData(data)
	}
	return
}

// RUN CODE =========================================================

// Subclasses (MockUpaxClient, etc) use sequences of calls to these
// functions to accomplish their purposes.

func (upc *UpaxClient) SessionSetup(proposedVersion uint32) (
	upcx *xt.TcpConnection, decidedVersion uint32, err error) {
	var (
		ciphertext1, key1, salt1, salt1c []byte
		ciphertext2, key2, salt2         []byte
	)
	// Set up connection to server. -----------------------------
	ctor, err := xt.NewTcpConnector(upc.serverEnd)
	if err == nil {
		var conn xt.ConnectionI
		conn, err = ctor.Connect(nil)
		if err == nil {
			upcx = conn.(*xt.TcpConnection)
		}
	}
	// Send HELLO -----------------------------------------------
	if err == nil {
		upc.Cnx = upcx
		ciphertext1, key1, salt1,
			err = xa.ClientEncryptHello(proposedVersion, upc.serverCK)
	}
	if err == nil {
		err = upc.WriteData(ciphertext1)
	}
	// Process HELLO REPLY --------------------------------------
	if err == nil {
		ciphertext2, err = upc.ReadData()
	}
	if err == nil {
		key2, salt2, salt1c, decidedVersion,
			err = xa.ClientDecryptHelloReply(ciphertext2, key1)
		_ = salt1c // XXX
	}
	// Set up AES engines ---------------------------------------
	if err == nil {
		upc.salt1 = salt1
		upc.key2 = key2
		upc.salt2 = salt2
		upc.Version = xu.DecimalVersion(decidedVersion)
		upc.engine, err = aes.NewCipher(key2)
		//if err == nil {
		//	upc.encrypter = cipher.NewCBCEncrypter(upc.engine, iv2)
		//	upc.decrypter = cipher.NewCBCDecrypter(upc.engine, iv2)
		//}
	}
	return
}

// msgN, token including DigSig; gets Ack or Error
func (upc *UpaxClient) IntroAndAck() (err error) {
	var (
		name                       string
		id, ckBytes, skBytes, salt []byte
		digSig                     []byte // over name, id, ckBytes, skBytes, salt, in order
	)
	// Send INTRO MSG =====================================
	name = upc.GetName()
	id = upc.GetNodeID().Value()
	ckBytes, err = xc.RSAPubKeyToWire(&upc.ckPriv.PublicKey)
	if err == nil {
		skBytes, err = xc.RSAPubKeyToWire(&upc.skPriv.PublicKey)
		if err == nil {
			rng := xr.NewSystemRNG(0)
			n := uint64(rng.Int63())
			salt = make([]byte, 8)
			binary.LittleEndian.PutUint64(salt, n)
		}
	}
	if err == nil {
		d := sha1.New()
		d.Write([]byte(name))
		d.Write(id)
		d.Write(ckBytes)
		d.Write(skBytes)
		d.Write(salt)
		hash := d.Sum(nil)
		digSig, err = rsa.SignPKCS1v15(
			rand.Reader, upc.skPriv, crypto.SHA1, hash)
	}
	if err == nil {
		token := UpaxClientMsg_Token{
			Name:     &name,
			ID:       id,
			CommsKey: ckBytes,
			SigKey:   skBytes,
			Salt:     salt,
			DigSig:   digSig,
		}
		op := UpaxClientMsg_Intro
		request := &UpaxClientMsg{
			Op:         &op,
			ClientInfo: &token,
		}
		// SHOULD CHECK FOR TIMEOUT
		err = upc.writeMsg(request)
	}
	// Process ACK ========================================
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S AN ACK
		response, err = upc.readMsg()
		op := response.GetOp()
		if op != UpaxClientMsg_Ack {
			err = ExpectedAck
		}
	}
	return
}

// msgN, id, salt, sig; gets Ack or Error
func (upc *UpaxClient) ItsMeAndAck() (err error) {
	var (
		id, salt []byte
		digSig   []byte // over id, salt, in order
	)
	// Send ITS_ME MSG ====================================

	id = upc.GetNodeID().Value()
	rng := xr.NewSystemRNG(0)
	n := uint64(rng.Int63())
	salt = make([]byte, 8)
	binary.LittleEndian.PutUint64(salt, n)

	d := sha1.New()
	d.Write(id)
	d.Write(salt)
	hash := d.Sum(nil)
	digSig, err = rsa.SignPKCS1v15(
		rand.Reader, upc.skPriv, crypto.SHA1, hash)

	if err == nil {
		op := UpaxClientMsg_ItsMe
		request := &UpaxClientMsg{
			Op:   &op,
			ID:   id,
			Salt: salt,
			Sig:  digSig,
		}
		// SHOULD CHECK FOR TIMEOUT
		err = upc.writeMsg(request)
	}
	// Process ACK ========================================
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT
		response, err = upc.readMsg()
		op := response.GetOp()
		if op != UpaxClientMsg_Ack {
			err = ExpectedAck
		}
	}

	return
}

// msgN ; gets Ack or timeout
func (upc *UpaxClient) KeepAliveAndAck() (err error) {

	// Send KEEP_ALIVE MSG ================================
	op := UpaxClientMsg_KeepAlive
	request := &UpaxClientMsg{
		Op: &op,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = upc.writeMsg(request)

	// Process ACK ========================================
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT
		response, err = upc.readMsg()
		op := response.GetOp()
		if op != UpaxClientMsg_Ack {
			err = ExpectedAck

			// XXX MsgN, YourMsgN ignored
		}
	}

	return
}

// msgN, hash ; gets Ack or NotFound -- LEAVE THIS STUBBED
func (upc *UpaxClient) QueryAndAck() (err error) {

	panic("QueryAndAck not implemented")

	// Process ACK ========================================
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT
		response, err = upc.readMsg()
		op := response.GetOp()
		if op != UpaxClientMsg_Ack {
			err = ExpectedAck

			// XXX MsgN, YourMsgN ignored
		}
	}
	return
}

// msgN, hash ; gets Data or NotFound
func (upc *UpaxClient) GetAndData(hash []byte) (
	logEntry *LogEntry, payload []byte, err error) {

	if hash == nil {
		err = NilHash
	} else {
		op := UpaxClientMsg_Get
		request := &UpaxClientMsg{
			Op:   &op,
			Hash: hash,
		}
		// SHOULD CHECK FOR TIMEOUT
		err = upc.writeMsg(request)
	}
	// Process DATA or NOT_FOUND ==========================
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT
		response, err = upc.readMsg()
		if err == nil {
			op := response.GetOp()
			if op == UpaxClientMsg_Data {
				// Data has msgN, LogEntry, and payload fields
				payload = response.GetPayload()
				entryMsg := response.GetEntry()
				if entryMsg == nil {
					// SUITABLE ERROR
				} else {
					// any field may be bad
					t := entryMsg.GetTimestamp()
					key := entryMsg.GetContentKey()
					owner := entryMsg.GetOwner() // []byte
					src := entryMsg.GetSrc()
					path := entryMsg.GetPath()
					logEntry, err = NewLogEntry(t, key, owner, src, path)
				}
			} else if op == UpaxClientMsg_NotFound {
				// NotFound has only msgN, yourMsgN fields; return nil logEntry
				// and payload

				// XXX not an error

			} else {
				// XXX STUB: INVALID OP
			}
		}
	}
	return
}

// msgN plus IHaves ; gets Ack -- LEAVE THIS STUBBED
func (upc *UpaxClient) IHaveAndAck() (err error) {

	panic("IHaveAndAck not implemented")

	// Process ACK ========================================
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT
		response, err = upc.readMsg()
		op := response.GetOp()
		if op != UpaxClientMsg_Ack {
			err = ExpectedAck

			// XXX MsgN, YourMsgN ignored
		}
	}
	return
}

// msgN, logEntry, payload ; gets Ack
func (upc *UpaxClient) PutAndAck(entry *LogEntry, payload []byte) (err error) {

	if entry == nil {
		err = NilLogEntry
	} else if payload == nil {
		err = NilPayload
	}
	if err == nil {
		t := entry.Timestamp()
		key := entry.Key()
		owner := entry.NodeID()
		src := entry.Src()
		path := entry.Path()
		entry := &UpaxClientMsg_LogEntry{
			Timestamp:  &t,
			ContentKey: key,
			Owner:      owner,
			Src:        &src,
			Path:       &path,
		}
		op := UpaxClientMsg_Put
		request := &UpaxClientMsg{
			Op:      &op,
			Entry:   entry,
			Payload: payload,
		}
		// SHOULD CHECK FOR TIMEOUT
		err = upc.writeMsg(request)
	}

	// Process ACK ========================================
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT
		response, err = upc.readMsg()
		op := response.GetOp()
		if op != UpaxClientMsg_Ack {
			err = ExpectedAck

			// XXX MsgN, YourMsgN ignored
		}
	}
	return
}

// Send Bye, wait for and process Ack.

func (upc *UpaxClient) ByeAndAck() (err error) {

	op := UpaxClientMsg_Bye
	request := &UpaxClientMsg{
		Op: &op,
	}
	// SHOULD CHECK FOR TIMEOUT
	err = upc.writeMsg(request)

	// Process ACK = BYE REPLY ----------------------------------
	if err == nil {
		var response *UpaxClientMsg

		// SHOULD CHECK FOR TIMEOUT AND VERIFY THAT IT'S AN ACK
		response, err = upc.readMsg()
		op := response.GetOp()
		_ = op
	}
	return
}
