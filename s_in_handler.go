package upax_go

// xlattice_go/upax_go/s_in_handler.go

import (
	"fmt"
	xcl "github.com/jddixon/xlCluster_go"
	xa "github.com/jddixon/xlProtocol_go/aes_cnx"
	reg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	xu "github.com/jddixon/xlU_go"
)

var _ = fmt.Print

// See s_states.go

const (
	// the number of valid states upon receiving a message from a peer
	S_IN_STATE_COUNT = S_BYE_RCVD + 1

	// The tags that ClusterInHandler will accept from a peer.
	S_MIN_TAG = uint(UpaxClusterMsg_ItsMe)
	S_MAX_TAG = uint(UpaxClusterMsg_Bye)

	S_MSG_HANDLER_COUNT = S_MAX_TAG + 1
)

var (
	sMsgHandlers [][]interface{}
)

type ClusterInHandler struct {
	us       *UpaxServer
	uDir     xu.UI
	peerInfo *xcl.MemberInfo

	myMsgN   uint64 // first message 1, then increment on each send
	peerMsgN uint64 // expect this to be 1 on the first message

	version    uint32 // protocol version used in session
	entryState int
	exitState  int
	msgIn      *UpaxClusterMsg
	msgOut     *UpaxClusterMsg
	errOut     error

	ClusterCnxHandler
}

// Given an open new connection, create a handler for the connection,
// associating the connection with a registry.
func NewClusterInHandler(us *UpaxServer, conn xt.ConnectionI) (
	h *ClusterInHandler, err error) {

	if us == nil {
		err = NilServer
	} else if us.uDir == nil {
		err = NilUDir
	} else if conn == nil {
		err = xa.NilConnection
	} else {
		cnx := conn.(*xt.TcpConnection)
		h = &ClusterInHandler{
			us:   us,
			uDir: us.uDir,
			ClusterCnxHandler: ClusterCnxHandler{
				Cnx: cnx,
			},
		}
	}
	return
}

// Convert a protobuf op into a zero-based tag for use in the
// ClusterInHandler's dispatch table.
func clusterOp2tag(op UpaxClusterMsg_Tag) uint {
	return uint(op - UpaxClusterMsg_ItsMe)
}

// Given a handler associating an open new connection with a registry,
// process a hello message for this node, which creates a session.
// The hello message contains an AES Key+IV, a salt, and a requested
// protocol version. The salt must be at least eight bytes long.
func (h *ClusterInHandler) Run() (err error) {

	defer func() {
		if h.Cnx != nil {
			h.Cnx.Close()
		}
	}()

	// This adds an AES iv2 and key2 to the handler.
	err = handleClusterHello(h)
	if err == nil {
		// Given iv2, key2 create encrypt and decrypt engines.
		err = h.SetupSessionKey()
	}
	for err == nil {
		var (
			tag uint
		)
		// REQUEST --------------------------------------------------
		//   receive the raw data off the wire
		var ciphertext []byte
		ciphertext, err = h.ReadData()
		if err == nil {
			h.msgIn, err = clusterDecryptUnpadDecode(ciphertext, h.decrypter)
		}
		if err != nil {
			break
		}
		op := h.msgIn.GetOp()
		// TODO: range check on either op or tag
		tag = clusterOp2tag(op)
		if tag < S_MIN_TAG || tag > S_MAX_TAG {
			h.errOut = reg.TagOutOfRange
		}
		// ACTION ----------------------------------------------------
		// Take the action appropriate for the current state
		sMsgHandlers[h.entryState][tag].(func(*ClusterInHandler))(h)

		// RESPONSE -------------------------------------------------
		// Convert any error encountered into an error message to be
		// sent to the client.
		if h.errOut != nil {
			h.us.Logger.Printf("errOut to client: %s\n", h.errOut.Error())

			op := UpaxClusterMsg_Error
			s := h.errOut.Error()
			h.msgOut = &UpaxClusterMsg{
				Op:      &op,
				ErrDesc: &s,
			}
			h.errOut = nil            // reduce potential for confusion
			h.exitState = S_IN_CLOSED // there is no recovery from errors
		}

		// encode, pad, and encrypt the UpaxClusterMsg object
		if h.msgOut != nil {
			ciphertext, err = clusterEncodePadEncrypt(h.msgOut, h.encrypter)

			// XXX log any error
			if err != nil {
				h.us.Logger.Printf(
					"ClusterInHandler.Run: clusterEncodePadEncrypt returns %s\n",
					err.Error())
			}

			// put the ciphertext on the wire
			if err == nil {
				err = h.WriteData(ciphertext)

				// log any error
				if err != nil {
					h.us.Logger.Printf(
						"ClusterInHandler.Run: WriteData returns %s\n",
						err.Error())
				}
			}
		}
		h.entryState = h.exitState
		if h.exitState == S_IN_CLOSED {
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
func handleClusterHello(h *ClusterInHandler) (err error) {
	var (
		ciphertext, iv1, key1, salt1 []byte
		version1                     uint32
	)
	ciphertext, err = h.ReadData()
	if err == nil {
		iv1, key1, salt1, version1,
			err = xa.ServerDecodeHello(ciphertext, h.us.ckPriv)
		_ = version1 // ignore whatever version they propose
	}
	if err == nil {
		version2 := serverVersion
		iv2, key2, salt2, ciphertextOut, err := xa.ServerEncodeHelloReply(
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
			h.State = S_HELLO_RCVD
		}
	}
	// On any error silently close the connection.
	if err != nil {
		// DEBUG
		fmt.Printf("handleClusterHello closing cnx, error was %s\n",
			err.Error())
		// END
		h.Cnx.Close()
	}
	return
}
