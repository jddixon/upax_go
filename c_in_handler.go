package upax_go

// xlattice_go/upax_go/c_in_handler.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xa "github.com/jddixon/xlProtocol_go/aes_cnx"
	xcl "github.com/jddixon/xlCluster_go"
	reg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	xu "github.com/jddixon/xlU_go"
)

var _ = fmt.Print

// See c_states.go

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
)

type ClientInHandler struct {
	us         *UpaxServer
	uDir       xu.UI
	peerInfo   *xcl.MemberInfo

	myMsgN   uint64 // first message 1, then increment on each send
	peerMsgN uint64 // expect this to be 1 on the first message

	version    uint32 // protocol version used in session
	entryState int
	exitState  int
	msgIn      *UpaxClientMsg
	msgOut     *UpaxClientMsg
	errOut	   error

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
		err = xa.NilConnection
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

// Convert a protobuf op into a zero-based tag for use in the 
// ClientInHandler's dispatch table.
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

	// This adds an AES key2 to the handler.
	err = handleClientHello(h)
	for err == nil {
		var (
			tag uint
		)
		// REQUEST --------------------------------------------------
		//   receive the raw data off the wire
		var ciphertext []byte
		ciphertext, err = h.ReadData()
		if err == nil {
			h.msgIn, err = clientDecryptUnpadDecode(
										ciphertext, h.decrypter)
		}
		if err != nil {
			break
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
			h.us.Logger.Printf("errOut to client: %s\n", h.errOut.Error())

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
			ciphertext, err = clientEncodePadEncrypt(h.msgOut, h.encrypter)

			// XXX log any error
			if err != nil {
				h.us.Logger.Printf(
					"ClientInHandler.Run: clientEncodePadEncrypt returns %s\n", 
					err.Error())
			}

			// put the ciphertext on the wire
			if err == nil {
				err = h.WriteData(ciphertext)

				// log any error
				if err != nil {
					h.us.Logger.Printf(
						"ClientInHandler.Run: WriteData returns %s\n", 
						err.Error())
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

/////////////////////////////////////////////////////////////////////
// XXX COMPARE WITH reg/ClientNode.SessionSetup, which is a better
// model for this code.
/////////////////////////////////////////////////////////////////////

// The client has sent the server a one-time AES key+iv encrypted with
// the server's RSA comms public key.  The server creates the real
// session iv+key and returns them to the client encrypted with the
// one-time key+iv.
func handleClientHello(h *ClientInHandler) (err error) {
	var (
		ciphertext, ciphertextOut []byte
		version1, version2                uint32
		sOneShot, sSession	*xa.AesSession
	)
	rng := xr.MakeSystemRNG()
	ciphertext, err = h.ReadData()
	if err == nil {
		sOneShot, version1,	err = xa.ServerDecryptHello(
			ciphertext, h.us.ckPriv, rng)
		_ = version1		// we don't actually use this
	}
	if err == nil {
		version2 = uint32(serverVersion)		// a global !
		sSession, ciphertextOut, err = xa.ServerEncryptHelloReply(
			sOneShot, version2)
		if err == nil {
			h.AesSession = *sSession
			err = h.WriteData(ciphertextOut)
		}
		if err == nil {
			h.version = version2
			h.State = C_HELLO_RCVD
		}
	}
	// On any error silently close the connection.
	if err != nil {
		// DEBUG
		fmt.Printf("handleClientHello closing cnx, error was %s\n", 
			err.Error())
		// END
		h.Cnx.Close()
	}
	return
}
