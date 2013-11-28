package upax_go

// upax_go/s_packets.go

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/aes"
	"crypto/cipher"
	xc "github.com/jddixon/xlattice_go/crypto"
	xt "github.com/jddixon/xlattice_go/transport"
)

const (
	S_MSG_BUF_LEN = 128 * 1024
)

type ClusterCnxHandler struct {
	State int
	Cnx   *xt.TcpConnection
}

// Read data from the connection.  
// XXX This will not handle partial reads correctly
func (h *ClusterCnxHandler) ReadData() (data []byte, err error) {
	data = make([]byte, S_MSG_BUF_LEN)
	count, err := h.Cnx.Read(data)
	if err == nil && count > 0 {
		data = data[:count]
		return
	}
	return nil, err
}

// Write data to the connection.
func (h *ClusterCnxHandler) WriteData(data []byte) (err error) {
	count, err := h.Cnx.Write(data)
	// XXX handle cases where not all bytes written
	_ = count
	return
}
func decodeClusterPacket(data []byte) (*UpaxClusterMsg, error) {
	var m UpaxClusterMsg
	err := proto.Unmarshal(data, &m)
	// XXX do some filtering, eg for nil op
	return &m, err
}

func encodeClusterPacket(msg *UpaxClusterMsg) (
	data []byte, err error) {

	return proto.Marshal(msg)
}

func clusterEncodePadEncrypt(msg *UpaxClusterMsg, engine cipher.BlockMode) (
	ciphertext []byte, err error) {

	var paddedData []byte
	cData, err := encodeClusterPacket(msg)
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

func clusterDecryptUnpadDecode(ciphertext []byte, engine cipher.BlockMode) (
	msg *UpaxClusterMsg, err error) {

	plaintext := make([]byte, len(ciphertext))
	engine.CryptBlocks(plaintext, ciphertext) // dest <- src

	unpaddedCData, err := xc.StripPKCS7Padding(plaintext, aes.BlockSize)
	if err == nil {
		msg, err = decodeClusterPacket(unpaddedCData)
	}
	return
}
