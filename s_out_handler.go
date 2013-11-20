package upax_go

// upax_go/s_out_handler.go

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"github.com/jddixon/xlattice_go/msg"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	"github.com/jddixon/xlattice_go/u"
	xu "github.com/jddixon/xlattice_go/util"
	"sync"
)

func init() {
	var err error
	serverVersion, err = xu.ParseDecimalVersion(VERSION)
	if err != nil {
		panic(err)
	}
}

// Upax servers open a connection to each other server in the cluster.
//
type ClusterOutHandler struct {
	iv1, key1, iv2, key2, salt1, salt2 []byte
	engineS                            cipher.Block
	encrypterS                         cipher.BlockMode
	decrypterS                         cipher.BlockMode

	us       *UpaxServer
	uDir     u.UI
	peerInfo *reg.MemberInfo

	myMsgN   uint64 // first message 1, then increment on each send
	peerMsgN uint64 // expect this to be 1 on the first message
	msgNMu   sync.Mutex

	version uint32 // protocol version used in session
	errOut  error
	ClusterCnxHandler
}

func NewClusterOutHandler(us *UpaxServer, conn xt.ConnectionI) (
	h *ClusterOutHandler, err error) {

	if us == nil {
		err = NilServer
	} else if us.uDir == nil {
		err = NilUDir
	} else if conn == nil {
		err = msg.NilConnection
	} else {
		cnx := conn.(*xt.TcpConnection)
		h = &ClusterOutHandler{
			us:   us,
			uDir: us.uDir,
			ClusterCnxHandler: ClusterCnxHandler{
				Cnx: cnx,
			},
		}
	}
	return
}

// XXX This duplicates ClusterInHandler.SetUpPeerSessionKey: some
// refactoring is needed.
//
func SetUpPeerSessionOutKey(h *ClusterOutHandler) (err error) {
	h.engineS, err = aes.NewCipher(h.key2)
	if err == nil {
		h.encrypterS = cipher.NewCBCEncrypter(h.engineS, h.iv2)
		h.decrypterS = cipher.NewCBCDecrypter(h.engineS, h.iv2)
	}
	return
}

// This is unnecessary if only one thread steps the message number.
//
func (h *ClusterOutHandler) stepMsgN() {
	h.msgNMu.Lock()
	h.myMsgN++
	h.msgNMu.Unlock()
}

// Given a handler associating an open new connection with a peer,
// process a hello message for this node, which creates a session.
// The hello message contains an AES Key+IV, a salt, and a requested
// protocol version. The salt must be at least eight bytes long.

func (h *ClusterOutHandler) Run() (err error) {

	defer func() {
		if h.Cnx != nil {
			h.Cnx.Close()
		}
	}()

	// This adds an AES iv2 and key2 to the handler.
	err = handleOutPeerHello(h)
	if err == nil {
		// Given iv2, key2 create encrypt and decrypt engines.
		err = SetUpPeerSessionOutKey(h)
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
//
// XXX This is simply a copy of the function in s_in_handler.go, with
// "Out" inserted: we definitely need to refactor!
//
func handleOutPeerHello(h *ClusterOutHandler) (err error) {
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
		fmt.Printf("handleOutPeerHello closing cnx, error was %v\n", err)
		// END
		h.Cnx.Close()
		h = nil
	}
	return
}
