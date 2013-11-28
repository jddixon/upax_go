package upax_go

// xlattice_go/upax_go/s_in_handler.go

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
	// States through which a intra-cluster input cnx may pass
	S_HELLO_RCVD = iota

	// After the peer has sent a message containing its nodeID and a
	// digital signature, we can determine which peer we are speaking to.
	S_ID_VERIFIED

	// Once the connection has reached this state, no more messages
	// can be accepted.
	S_BYE_RCVD

	// When we reach this state, the connection must be closed.
	S_IN_CLOSED
)

const (
	// the number of valid states upon receiving a message from a peer
	S_IN_STATE_COUNT = S_BYE_RCVD + 1

	// The tags that ClusterInHandler will accept from a peer.
	S_MIN_TAG = uint(UpaxClusterMsg_ItsMe)
	S_MAX_TAG = uint(UpaxClusterMsg_Bye)

	S_MSG_HANDLER_COUNT = S_MAX_TAG + 1
)

var (
	sMsgHandlers  [][]interface{}
	serverVersion xu.DecimalVersion
)

func init() {
	// sMsgHandlers = make([][]interface{}, S_BYE_RCVD, S_MSG_HANDLER_COUNT)

	sMsgHandlers = [][]interface{}{
		// client messages permitted in S_HELLO_RCVD state
		{doSItsMeMsg, badSCombo, badSCombo, badSCombo, badSCombo, badSCombo},
		// client messages permitted in S_ID_VERIFIED state
		{badSCombo, doSKeepAliveMsg, doSGetMsg, doSIHaveMsg, doSPutMsg, doSByeMsg},
	}
	var err error
	serverVersion, err = xu.ParseDecimalVersion(VERSION)
	if err != nil {
		panic(err)
	}
}

type ClusterInHandler struct {
	us       *UpaxServer
	uDir     u.UI
	peerInfo *reg.MemberInfo

	myMsgN   uint64 // first message 1, then increment on each send
	peerMsgN uint64 // expect this to be 1 on the first message

	version    uint32 // protocol version used in session
	entryState int
	exitState  int
	msgIn      *UpaxClusterMsg
	msgOut     *UpaxClusterMsg
	errOut     error
	
	engineS                            cipher.Block
	encrypterS                         cipher.BlockMode
	decrypterS                         cipher.BlockMode
	iv1, key1, iv2, key2, salt1, salt2 []byte

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
		err = msg.NilConnection
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

func SetUpPeerSessionKey(h *ClusterInHandler) (err error) {
	h.engineS, err = aes.NewCipher(h.key2)
	if err == nil {
		h.encrypterS = cipher.NewCBCEncrypter(h.engineS, h.iv2)
		h.decrypterS = cipher.NewCBCDecrypter(h.engineS, h.iv2)
	}
	return
}

// Convert a protobuf op into a zero-based tag for use in the ClusterInHandler's
// dispatch table.
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
	err = handlePeerHello(h)
	if err != nil {
		return
	}
	// Given iv2, key2 create encrypt and decrypt engines.
	err = SetUpPeerSessionKey(h)
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
		h.msgIn, err = clusterDecryptUnpadDecode(ciphertext, h.decrypterS)
		if err != nil {
			return
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
			h.us.Logger.Printf("errOut to client: %v\n", h.errOut)

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
			ciphertext, err = clusterEncodePadEncrypt(h.msgOut, h.encrypterS)

			// XXX log any error
			if err != nil {
				h.us.Logger.Printf(
					"ClusterInHandler.Run: clusterEncodePadEncrypt returns %v\n", err)
			}

			// put the ciphertext on the wire
			if err == nil {
				err = h.WriteData(ciphertext)

				// XXX log any error
				if err != nil {
					h.us.Logger.Printf(
						"ClusterInHandler.Run: WriteData returns %v\n", err)
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

// The client has sent the server a one-time AES key+iv encrypted with
// the server's RSA comms public key.  The server creates the real
// session iv+key and returns them to the client encrypted with the
// one-time key+iv.


/////////////////////////////////////////////////////////////////////
// XXX THIS IS WRONG.  COMPARE WITH reg/ClientNode.SessionSetup, which 
// is the right model for this code.
/////////////////////////////////////////////////////////////////////
func handlePeerHello(h *ClusterInHandler) (err error) {
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
			h.State = S_HELLO_RCVD
		}
	}
	// On any error silently close the connection and delete the handler,
	// an exciting thing to do.
	if err != nil {
		// DEBUG
		fmt.Printf("handlePeerHello closing cnx, error was %v\n", err)
		// END
		h.Cnx.Close()
		h = nil
	}
	return
}
