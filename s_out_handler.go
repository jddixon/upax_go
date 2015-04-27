package upax_go

// upax_go/s_out_handler.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xcl "github.com/jddixon/xlCluster_go"
	xa "github.com/jddixon/xlProtocol_go/aes_cnx"
	xt "github.com/jddixon/xlTransport_go"
	u "github.com/jddixon/xlU_go"
	xu "github.com/jddixon/xlUtil_go"
	"sync"
)

const (
	OUT_MSG_Q_SIZE = 16
)

func init() {
	var err error
	serverVersion, err = xu.ParseDecimalVersion(VERSION)
	if err != nil {
		panic(err)
	}
}

// Each Upax server opens a connection to each other server in the cluster.
//
type ClusterOutHandler struct {
	//key1, key2, salt1, salt2 []byte
	//engineS                  cipher.Block
	//encrypterS               cipher.BlockMode
	//decrypterS               cipher.BlockMode

	us       *UpaxServer
	uDir     u.UI
	peerInfo *xcl.MemberInfo

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
		err = xa.NilConnection
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

// XXX This duplicates ClusterInHandler.SetupPeerSessionKey: some
// refactoring is needed.
//
//func SetupPeerSessionOutKey(h *ClusterOutHandler) (err error) {
//	h.engineS, err = aes.NewCipher(h.key2)
//	if err == nil {
//		h.encrypterS = cipher.NewCBCEncrypter(h.engineS, h.iv2)
//		h.decrypterS = cipher.NewCBCDecrypter(h.engineS, h.iv2)
//	}
//	return
//}

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

	// This adds an sSession (AES key2) to the handler.
	err = handleOutPeerHello(h)
	//if err == nil {
	//	// Given iv2, key2 create encrypt and decrypt engines.
	//	err = SetupPeerSessionOutKey(h)
	//}
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
		ciphertext, ciphertextOut []byte
		version1                  uint32
		sOneShot, sSession        *xa.AesSession
		rng                       *xr.PRNG
	)
	ciphertext, err = h.ReadData()
	if err == nil {
		rng = xr.MakeSystemRNG()
		sOneShot, version1, err = xa.ServerDecryptHello(
			ciphertext, h.us.ckPriv, rng)
	}
	if err == nil {
		_ = version1 // just ignored for now
		version2 := uint32(serverVersion)
		sSession, ciphertextOut, err = xa.ServerEncryptHelloReply(
			sOneShot, version2)
		if err == nil {
			h.AesSession = *sSession
			err = h.WriteData(ciphertextOut)
		}
		if err == nil {
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
