package upax_go

// upax_go/c_packets.go

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xt "github.com/jddixon/xlattice_go/transport"
	// "sync"
)

var _ = fmt.Print

const (
	C_MSG_BUF_LEN = 128 * 1024
)

type ClientCnxHandler struct {
	State int
	Cnx   *xt.TcpConnection
}

// Read data from the connection.  XXX This will not handle partial
// reads correctly
//
func (h *ClientCnxHandler) ReadData() (data []byte, err error) {
	data = make([]byte, C_MSG_BUF_LEN)
	count, err := h.Cnx.Read(data)
	// DEBUG
	//fmt.Printf("ReadData: count is %d, err is %v\n", count, err)
	// END
	if err == nil && count > 0 {
		data = data[:count]
		return
	}
	return nil, err
}

func (h *ClientCnxHandler) WriteData(data []byte) (err error) {
	count, err := h.Cnx.Write(data)
	// XXX handle cases where not all bytes written
	_ = count
	return
}
func decodeClientPacket(data []byte) (*UpaxClientMsg, error) {
	var m UpaxClientMsg
	err := proto.Unmarshal(data, &m)
	// XXX do some filtering, eg for nil op
	return &m, err
}

func encodeClientPacket(msg *UpaxClientMsg) (data []byte, err error) {
	return proto.Marshal(msg)
}

func clientEncodePadEncrypt(msg *UpaxClientMsg, engine cipher.BlockMode) (
	ciphertext []byte, err error) {

	var paddedData []byte
	cData, err := encodeClientPacket(msg)
	if err == nil {
		paddedData, err = xc.AddPKCS7Padding(cData, aes.BlockSize)
	}
	if err == nil {
		msgLen := len(paddedData)
		nBlocks := (msgLen + aes.BlockSize - 2) / aes.BlockSize
		ciphertext = make([]byte, nBlocks*aes.BlockSize)
		engine.CryptBlocks(ciphertext, paddedData) // dest <- src
	}
	return
}

func clientDecryptUnpadDecode(ciphertext []byte, engine cipher.BlockMode) (
	msg *UpaxClientMsg, err error) {

	plaintext := make([]byte, len(ciphertext))
	engine.CryptBlocks(plaintext, ciphertext) // dest <- src

	unpaddedCData, err := xc.StripPKCS7Padding(plaintext, aes.BlockSize)
	if err == nil {
		msg, err = decodeClientPacket(unpaddedCData)
	}
	return
}
