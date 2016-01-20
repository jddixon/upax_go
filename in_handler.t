package upax_go

// xlattice_go/${pkgName}/${filePrefix}in_handler.go

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
	uDir       xu.UI
	peerInfo   *xcl.MemberInfo

	myMsgN   uint64 // first message 1, then increment on each send
	peerMsgN uint64 // expect this to be 1 on the first message

	version    uint32 // protocol version used in session
	entryState int
	exitState  int
	msgIn      *Upax${TypePrefix}Msg
	msgOut     *Upax${TypePrefix}Msg
	errOut	   error

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
		err = xa.NilConnection
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

// Convert a protobuf op into a zero-based tag for use in the 
// ${TypePrefix}InHandler's dispatch table.
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

	// This adds an AES key2 to the handler.
	err = handle${TypePrefix}Hello(h)
	for err == nil {
		var (
			tag uint
		)
		// REQUEST --------------------------------------------------
		//   receive the raw data off the wire
		var ciphertext []byte
		ciphertext, err = h.ReadData()
		if err == nil {
			h.msgIn, err = ${funcPrefix}DecryptUnpadDecode(
										ciphertext, h.decrypter)
		}
		if err != nil {
			break
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
			h.us.Logger.Printf("errOut to client: %s\n", h.errOut.Error())

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
			ciphertext, err = ${funcPrefix}EncodePadEncrypt(h.msgOut, h.encrypter)

			// XXX log any error
			if err != nil {
				h.us.Logger.Printf(
					"${TypePrefix}InHandler.Run: ${funcPrefix}EncodePadEncrypt returns %s\n", 
					err.Error())
			}

			// put the ciphertext on the wire
			if err == nil {
				err = h.WriteData(ciphertext)

				// log any error
				if err != nil {
					h.us.Logger.Printf(
						"${TypePrefix}InHandler.Run: WriteData returns %s\n", 
						err.Error())
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

/////////////////////////////////////////////////////////////////////
// XXX COMPARE WITH reg/ClientNode.SessionSetup, which is a better
// model for this code.
/////////////////////////////////////////////////////////////////////

// The client has sent the server a one-time AES key+iv encrypted with
// the server's RSA comms public key.  The server creates the real
// session iv+key and returns them to the client encrypted with the
// one-time key+iv.
func handle${TypePrefix}Hello(h *${TypePrefix}InHandler) (err error) {
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
			h.State = ${ConstPrefix}HELLO_RCVD
		}
	}
	// On any error silently close the connection.
	if err != nil {
		// DEBUG
		fmt.Printf("handle${TypePrefix}Hello closing cnx, error was %s\n", 
			err.Error())
		// END
		h.Cnx.Close()
	}
	return
}
