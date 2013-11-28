package upax_go

// xlattice_go/upax_go/c_in_handler.go

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/jddixon/xlattice_go/msg"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	"github.com/jddixon/xlattice_go/u"
	xu "github.com/jddixon/xlattice_go/util"
)

var _ = fmt.Print

// For server-server connections, there is little state to track.
const (
	// States through which client-server input cnx may pass
	C_HELLO_RCVD = iota

	// After the peer has sent a message containing its nodeID and a
	// digital signature, we can determine which peer we are speaking to.
	C_ID_VERIFIED

	// Once the connection has reached this state, no more messages
	// can be accepted.
	C_BYE_RCVD

	// When we reach this state, the connection must be closed.
	C_IN_CLOSED
)

const (
	// the number of valid states upon receiving a message from a peer
	C_IN_STATE_COUNT = C_BYE_RCVD + 1

	// The tags that ClientInHandler will accept from a peer.
	C_MIN_TAG = uint(UpaxClientMsg_Intro)
	C_MAX_TAG = uint(UpaxClientMsg_Bye)

	C_MSG_HANDLER_COUNT = C_MAX_TAG + 1
)

var (
	cMsgHandlers [][]interface{}
	// serverVersion xu.DecimalVersion	// declared in s_in_handler.go
)

func init() {
	// cMsgHandlers = make([][]interface{}, C_BYE_RCVD, C_MSG_HANDLER_COUNT)

	cMsgHandlers = [][]interface{}{
		// client messages permitted in C_HELLO_RCVD state
		{doCIntroMsg, doCItsMeMsg, badCCombo, badCCombo, badCCombo,
			badCCombo, badCCombo},
		// client messages permitted in C_ID_VERIFIED state
		{badCCombo, badCCombo, doCKeepAliveMsg, doCQueryMsg, doCGetMsg,
			doCPutMsg, doCByeMsg},
	}
	var err error
	serverVersion, err = xu.ParseDecimalVersion(VERSION)
	if err != nil {
		panic(err)
	}
}

type ClientInHandler struct {
	us         *UpaxServer
	uDir       u.UI
	clientInfo *reg.MemberInfo
	// client  *reg.RegClient

	myMsgN   uint64 // first message 1, then increment on each send
	peerMsgN uint64 // expect this to be 1 on the first message

	version    uint32 // protocol version used in session
	entryState int
	exitState  int
	msgIn      *UpaxClientMsg
	msgOut     *UpaxClientMsg
	errOut     error
	
	engineS                            cipher.Block
	encrypterS                         cipher.BlockMode
	decrypterS                         cipher.BlockMode
	iv1, key1, iv2, key2, salt1, salt2 []byte

	ClientCnxHandler
}

// Given an open new connection, create a handler for the connection,
// associating the connection with a registry.

func NewClientInHandler(us *UpaxServer, conn xt.ConnectionI) (
	h *ClientInHandler, err error) {

	if us == nil {
		err = NilServer
	} else if us.uDir == nil {
		err = NilUDir
	} else if conn == nil {
		err = msg.NilConnection
	} else {
		cnx := conn.(*xt.TcpConnection)
		h = &ClientInHandler{
			us:   us,
			uDir: us.uDir,
			ClientCnxHandler: ClientCnxHandler{
				Cnx: cnx,
			},
		}
	}
	return
}

// Set up the receiver (server) side of a communications link with 
// RSA-to-AES handshaking
//
func SetUpReceiverSessionKey(h *ClientInHandler) (err error) {
	h.engineS, err = aes.NewCipher(h.key2)
	if err == nil {
		h.encrypterS = cipher.NewCBCEncrypter(h.engineS, h.iv2)
		h.decrypterS = cipher.NewCBCDecrypter(h.engineS, h.iv2)
	}
	return
}

// Convert a protobuf op into a zero-based tag for use in the ClientInHandler's
// dispatch table.
func clientOp2tag(op UpaxClientMsg_Tag) uint {
	return uint(op - UpaxClientMsg_Intro)
}

// Given a handler associating an open new connection with a registry,
// process a hello message for this node, which creates a session.
// The hello message contains an AES Key+IV, a salt, and a requested
// protocol version. The salt must be at least eight bytes long.

func (h *ClientInHandler) Run() (err error) {

	defer func() {
		if h.Cnx != nil {
			h.Cnx.Close()
		}
	}()

	// This adds an AES iv2 and key2 to the handler.
	err = handleClientHello(h)
	if err != nil {
		return
	}
	// Given iv2, key2 create encrypt and decrypt engines.
	err = SetUpReceiverSessionKey(h)
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
		h.msgIn, err = clientDecryptUnpadDecode(ciphertext, h.decrypterS)
		if err != nil {
			return
		}
		op := h.msgIn.GetOp()
		// TODO: range check on either op or tag
		tag = clientOp2tag(op)
		if tag < C_MIN_TAG || tag > C_MAX_TAG {
			h.errOut = reg.TagOutOfRange
		}
		// ACTION ----------------------------------------------------
		// Take the action appropriate for the current state
		cMsgHandlers[h.entryState][tag].(func(*ClientInHandler))(h)

		// RESPONSE -------------------------------------------------
		// Convert any error encountered into an error message to be
		// sent to the client.
		if h.errOut != nil {
			h.us.Logger.Printf("errOut to client: %v\n", h.errOut)

			op := UpaxClientMsg_Error
			s := h.errOut.Error()
			h.msgOut = &UpaxClientMsg{
				Op:      &op,
				ErrDesc: &s,
			}
			h.errOut = nil            // reduce potential for confusion
			h.exitState = C_IN_CLOSED // there is no recovery from errors
		}

		// encode, pad, and encrypt the UpaxClientMsg object
		if h.msgOut != nil {
			ciphertext, err = clientEncodePadEncrypt(h.msgOut, h.encrypterS)

			// XXX log any error
			if err != nil {
				h.us.Logger.Printf(
					"ClientInHandler.Run: clientEncodePadEncrypt returns %v\n", err)
			}

			// put the ciphertext on the wire
			if err == nil {
				err = h.WriteData(ciphertext)

				// XXX log any error
				if err != nil {
					h.us.Logger.Printf(
						"ClientInHandler.Run: WriteData returns %v\n", err)
				}
			}

		}
		h.entryState = h.exitState
		if h.exitState == C_IN_CLOSED {
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

func handleClientHello(h *ClientInHandler) (err error) {
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
			h.State = C_HELLO_RCVD
		}
	}
	// On any error silently close the connection and delete the handler,
	// an exciting thing to do.
	if err != nil {
		// DEBUG
		fmt.Printf("handleClientHello closing cnx, error was %v\n", err)
		// END
		h.Cnx.Close()
		h = nil
	}
	return
}
