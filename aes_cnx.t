package ${pkgName} 

// ${pkgName}/${filePrefix}aes_cnx.go

import (
	"code.google.com/p/goprotobuf/proto"
	"crypto/aes"
	"crypto/cipher"
	xc "github.com/jddixon/xlattice_go/crypto"
	xt "github.com/jddixon/xlattice_go/transport"
)

const (
	${ConstPrefix}MSG_BUF_LEN = 128 * 1024
)

type ${TypePrefix}CnxHandler struct {
	State int
	Cnx   *xt.TcpConnection
	engine                            cipher.Block
	encrypter                         cipher.BlockMode
	decrypter                         cipher.BlockMode
	iv1, key1, iv2, key2, salt1, salt2 []byte
}

func (a *${TypePrefix}CnxHandler) SetupSessionKey() (err error) {
	a.engine, err = aes.NewCipher(a.key2)
	if err == nil {
		a.encrypter = cipher.NewCBCEncrypter(a.engine, a.iv2)
		a.decrypter = cipher.NewCBCDecrypter(a.engine, a.iv2)
	}
	return
}

// Read data from the connection.  
// XXX This will not handle partial reads correctly
func (h *${TypePrefix}CnxHandler) ReadData() (data []byte, err error) {
	data = make([]byte, ${ConstPrefix}MSG_BUF_LEN)
	count, err := h.Cnx.Read(data)
	if err == nil && count > 0 {
		data = data[:count]
		return
	}
	return nil, err
}

// Write data to the connection.
func (h *${TypePrefix}CnxHandler) WriteData(data []byte) (err error) {
	count, err := h.Cnx.Write(data)
	// XXX handle cases where not all bytes written
	_ = count
	return
}
func decode${TypePrefix}Packet(data []byte) (*Upax${TypePrefix}Msg, error) {
	var m Upax${TypePrefix}Msg
	err := proto.Unmarshal(data, &m)
	// XXX do some filtering, eg for nil op
	return &m, err
}

func encode${TypePrefix}Packet(msg *Upax${TypePrefix}Msg) (
	data []byte, err error) {

	return proto.Marshal(msg)
}

func ${funcPrefix}EncodePadEncrypt(msg *Upax${TypePrefix}Msg, engine cipher.BlockMode) (
	ciphertext []byte, err error) {

	var paddedData []byte
	cData, err := encode${TypePrefix}Packet(msg)
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

func ${funcPrefix}DecryptUnpadDecode(ciphertext []byte, engine cipher.BlockMode) (
	msg *Upax${TypePrefix}Msg, err error) {

	plaintext := make([]byte, len(ciphertext))
	engine.CryptBlocks(plaintext, ciphertext) // dest <- src

	unpaddedCData, err := xc.StripPKCS7Padding(plaintext, aes.BlockSize)
	if err == nil {
		msg, err = decode${TypePrefix}Packet(unpaddedCData)
	}
	return
}
