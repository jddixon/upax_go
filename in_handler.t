package ${pkgName}

// xlattice_go/${pkgName}/${filePrefix}in_handler.go

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/jddixon/xlattice_go/msg"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	"github.com/jddixon/xlattice_go/u"
)

var _ = fmt.Print

// See ${filePrefix}states.go

const (
	// the number of valid states upon receiving a message from a peer
	${ConstPrefix}IN_STATE_COUNT = ${ConstPrefix}BYE_RCVD + 1

	// The tags that ${TypePrefix}InHandler will accept from a peer.
	${ConstPrefix}MIN_TAG = uint(Upax${TypePrefix}Msg_${firstMsg})
	${ConstPrefix}MAX_TAG = uint(Upax${TypePrefix}Msg_Bye)

	${ConstPrefix}MSG_HANDLER_COUNT = ${ConstPrefix}MAX_TAG + 1
)

var (
	${shortPrefix}MsgHandlers [][]interface{}
)

type ${TypePrefix}InHandler struct {
	us         *UpaxServer
	uDir       u.UI
	peerInfo   *reg.MemberInfo
	// client  *reg.RegClient

	myMsgN   uint64 // first message 1, then increment on each send
	peerMsgN uint64 // expect this to be 1 on the first message

	version    uint32 // protocol version used in session
	entryState int
	exitState  int
	msgIn      *Upax${TypePrefix}Msg
	msgOut     *Upax${TypePrefix}Msg
	errOut     error
	
	engineS                            cipher.Block
	encrypterS                         cipher.BlockMode
	decrypterS                         cipher.BlockMode
	iv1, key1, iv2, key2, salt1, salt2 []byte

	${TypePrefix}CnxHandler
}

// Given an open new connection, create a handler for the connection,
// associating the connection with a registry.

func New${TypePrefix}InHandler(us *UpaxServer, conn xt.ConnectionI) (
	h *${TypePrefix}InHandler, err error) {

	if us == nil {
		err = NilServer
	} else if us.uDir == nil {
		err = NilUDir
	} else if conn == nil {
		err = msg.NilConnection
	} else {
		cnx := conn.(*xt.TcpConnection)
		h = &${TypePrefix}InHandler{
			us:   us,
			uDir: us.uDir,
			${TypePrefix}CnxHandler: ${TypePrefix}CnxHandler{
				Cnx: cnx,
			},
		}
	}
	return
}

// Set up the receiver (server) side of a communications link with 
// RSA-to-AES handshaking
//
func SetUp${TypePrefix}SessionKey(h *${TypePrefix}InHandler) (err error) {
	h.engineS, err = aes.NewCipher(h.key2)
	if err == nil {
		h.encrypterS = cipher.NewCBCEncrypter(h.engineS, h.iv2)
		h.decrypterS = cipher.NewCBCDecrypter(h.engineS, h.iv2)
	}
	return
}

// Convert a protobuf op into a zero-based tag for use in the ${TypePrefix}InHandler's
// dispatch table.
func ${funcPrefix}Op2tag(op Upax${TypePrefix}Msg_Tag) uint {
	return uint(op - Upax${TypePrefix}Msg_${firstMsg})
}

// Given a handler associating an open new connection with a registry,
// process a hello message for this node, which creates a session.
// The hello message contains an AES Key+IV, a salt, and a requested
// protocol version. The salt must be at least eight bytes long.

func (h *${TypePrefix}InHandler) Run() (err error) {

	defer func() {
		if h.Cnx != nil {
			h.Cnx.Close()
		}
	}()

	// This adds an AES iv2 and key2 to the handler.
	err = handle${TypePrefix}Hello(h)
	if err != nil {
		return
	}
	// Given iv2, key2 create encrypt and decrypt engines.
	err = SetUp${TypePrefix}SessionKey(h)
	if err != nil {
		return
	}
	for {
		var (
			tag uint
		)
		// REQUEST --------------------------------------------------
		//   receive the raw data off the wire
		var ciphertext []byte
		ciphertext, err = h.ReadData()
		if err != nil {
			return
		}
		h.msgIn, err = ${funcPrefix}DecryptUnpadDecode(ciphertext, h.decrypterS)
		if err != nil {
			return
		}
		op := h.msgIn.GetOp()
		// TODO: range check on either op or tag
		tag = ${funcPrefix}Op2tag(op)
		if tag < ${ConstPrefix}MIN_TAG || tag > ${ConstPrefix}MAX_TAG {
			h.errOut = reg.TagOutOfRange
		}
		// ACTION ----------------------------------------------------
		// Take the action appropriate for the current state
		${shortPrefix}MsgHandlers[h.entryState][tag].(func(*${TypePrefix}InHandler))(h)

		// RESPONSE -------------------------------------------------
		// Convert any error encountered into an error message to be
		// sent to the client.
		if h.errOut != nil {
			h.us.Logger.Printf("errOut to client: %v\n", h.errOut)

			op := Upax${TypePrefix}Msg_Error
			s := h.errOut.Error()
			h.msgOut = &Upax${TypePrefix}Msg{
				Op:      &op,
				ErrDesc: &s,
			}
			h.errOut = nil            // reduce potential for confusion
			h.exitState = ${ConstPrefix}IN_CLOSED // there is no recovery from errors
		}

		// encode, pad, and encrypt the Upax${TypePrefix}Msg object
		if h.msgOut != nil {
			ciphertext, err = ${funcPrefix}EncodePadEncrypt(h.msgOut, h.encrypterS)

			// XXX log any error
			if err != nil {
				h.us.Logger.Printf(
					"${TypePrefix}InHandler.Run: ${funcPrefix}EncodePadEncrypt returns %v\n", err)
			}

			// put the ciphertext on the wire
			if err == nil {
				err = h.WriteData(ciphertext)

				// XXX log any error
				if err != nil {
					h.us.Logger.Printf(
						"${TypePrefix}InHandler.Run: WriteData returns %v\n", err)
				}
			}

		}
		h.entryState = h.exitState
		if h.exitState == ${ConstPrefix}IN_CLOSED {
			break
		}
	}

	return
}

/////////////////////////////////////////////////////////////////////
// RSA-BASED MESSAGE PAIR
/////////////////////////////////////////////////////////////////////

// The client has sent the server a one-time AES key+iv encrypted with
// the server's RSA comms public key.  The server creates the real
// session iv+key and returns them to the client encrypted with the
// one-time key+iv.


/////////////////////////////////////////////////////////////////////
// XXX THIS IS WRONG.  COMPARE WITH reg/ClientNode.SessionSetup, which 
// is the right model for this code.
/////////////////////////////////////////////////////////////////////
func handle${TypePrefix}Hello(h *${TypePrefix}InHandler) (err error) {
	var (
		ciphertext, iv1, key1, salt1 []byte
		version1                     uint32
	)
	ciphertext, err = h.ReadData()
	if err == nil {
		iv1, key1, salt1, version1,
			err = msg.ServerDecodeHello(ciphertext, h.us.ckPriv)
		_ = version1 // ignore whatever version they propose
	}
	if err == nil {
		version2 := serverVersion
		iv2, key2, salt2, ciphertextOut, err := msg.ServerEncodeHelloReply(
			iv1, key1, salt1, uint32(version2))
		if err == nil {
			err = h.WriteData(ciphertextOut)
		}
		if err == nil {
			h.iv1 = iv1
			h.key1 = key1
			h.iv2 = iv2
			h.key2 = key2
			h.salt1 = salt1
			h.salt2 = salt2
			h.version = uint32(version2)
			h.State = ${ConstPrefix}HELLO_RCVD
		}
	}
	// On any error silently close the connection and delete the handler,
	// an exciting thing to do.
	if err != nil {
		// DEBUG
		fmt.Printf("handle${TypePrefix}Hello closing cnx, error was %v\n", err)
		// END
		h.Cnx.Close()
		h = nil
	}
	return
}