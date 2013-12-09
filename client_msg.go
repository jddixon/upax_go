package upax_go

// upax_go/client_msg.go

import (
	"crypto/aes"
	"crypto/cipher"

	// xc "github.com/jddixon/xlattice_go/crypto"
	xm "github.com/jddixon/xlattice_go/msg"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
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
		ciphertext1, iv1, key1, salt1, salt1c []byte
		ciphertext2, iv2, key2, salt2         []byte
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
		ciphertext1, iv1, key1, salt1,
			err = xm.ClientEncodeHello(proposedVersion, upc.serverCK)
	}
	if err == nil {
		err = upc.WriteData(ciphertext1)
	}
	// Process HELLO REPLY --------------------------------------
	if err == nil {
		ciphertext2, err = upc.ReadData()
	}
	if err == nil {
		iv2, key2, salt2, salt1c, decidedVersion,
			err = xm.ClientDecodeHelloReply(ciphertext2, iv1, key1)
		_ = salt1c // XXX
	}
	// Set up AES engines ---------------------------------------
	if err == nil {
		upc.salt1 = salt1
		upc.iv2 = iv2
		upc.key2 = key2
		upc.salt2 = salt2
		upc.Version = xu.DecimalVersion(decidedVersion)
		upc.engine, err = aes.NewCipher(key2)
		if err == nil {
			upc.encrypter = cipher.NewCBCEncrypter(upc.engine, iv2)
			upc.decrypter = cipher.NewCBCDecrypter(upc.engine, iv2)
		}
	}
	return
}

// msgN, token including DigSig; gets Ack or Error
func (upc *UpaxClient) IntroAndAck() (err error) {
	// XXX STUB XXX

	return
}

// msgN, id, opt salt, sig; gets Ack or Error
func (upc *UpaxClient) ItsMeAndAck() (err error) {
	// XXX STUB XXX

	return
}

// EXAMPLE: /////////////////////////////////////////////////////////
//
//func (upc *UpaxClient) ClientAndOK() (err error) {
//
//	var (
//		ckBytes, skBytes []byte
//		myEnds           []string
//	)
//
//	// Send CLIENT MSG ==========================================
//	ckBytes, err = xc.RSAPubKeyToWire(&upc.ckPriv.PublicKey)
//	if err == nil {
//		skBytes, err = xc.RSAPubKeyToWire(&upc.skPriv.PublicKey)
//		if err == nil {
//			for i := 0; i < len(upc.endPoints); i++ {
//				myEnds = append(myEnds, upc.endPoints[i].String())
//			}
//			token := &UpaxClientMsg_Token{
//				Name:     &upc.name,
//				Attrs:    &upc.proposedAttrs,
//				CommsKey: ckBytes,
//				SigKey:   skBytes,
//				MyEnds:   myEnds,
//			}
//
//			op := UpaxClientMsg_Client
//			request := &UpaxClientMsg{
//				Op:          &op,
//				ClientName:  &upc.name, // XXX redundant
//				ClientSpecs: token,
//			}
//			// SHOULD CHECK FOR TIMEOUT
//			err = upc.writeMsg(request)
//		}
//	}
//	// Process CLIENT_OK --------------------------------------------
//	// SHOULD CHECK FOR TIMEOUT
//	response, err := upc.readMsg()
//	if err == nil {
//		id := response.GetClientID()
//		upc.clientID, err = xi.New(id)
//
//		// XXX err ignored
//
//		upc.Attrs = response.GetClientAttrs()
//	}
//	return
//}

// msgN ; gets Ack or timeout
func (upc *UpaxClient) KeepAliveAndAck() (err error) {
	// XXX STUB XXX

	return
}

// msgN, hash ; gets Ack or NotFound -- LEAVE THIS STUBBED
func (upc *UpaxClient) QueryAndAck() (err error) {
	// XXX STUB XXX

	return
}

// msgN, hash ; gets Data or NotFound
func (upc *UpaxClient) GetAndData() (err error) {
	// XXX STUB XXX

	return
}

// msgN plus IHaves ; gets Ack -- LEAVE THIS STUBBED
func (upc *UpaxClient) IHaveAndAck() (err error) {
	// XXX STUB XXX

	return
}

// msgN, logEntry, payload ; gets Ack
func (upc *UpaxClient) PutAndAck() (err error) {
	// XXX STUB XXX

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
