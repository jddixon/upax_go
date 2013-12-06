package upax_go

import (
	"crypto/rand"
	"crypto/rsa"
	xn "github.com/jddixon/xlattice_go/node"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	xf "github.com/jddixon/xlattice_go/util/lfs"
)

// upax_go/mock_upax_client_test.go

type MockUpaxClient struct {
	K3     int // number of data items
	L1, L2 int // min and max length thereof
	UpaxClient
}

func NewMockUpaxClient(name, lfs string, members []*reg.MemberInfo) (
	mc *MockUpaxClient, err error) {

	var (
		ckPriv, skPriv *rsa.PrivateKey
		ep             []xt.EndPointI
		node           *xn.Node
		uc             *UpaxClient
	)

	// lfs should be a well-formed POSIX path; if the directory does
	// not exist we should create it.
	err = xf.CheckLFS(lfs)

	// The ckPriv is an RSA key used to encrypt short messages.
	if err == nil {
		if ckPriv == nil {
			ckPriv, err = rsa.GenerateKey(rand.Reader, 2048)
		}
		if err == nil {
			// The skPriv is an RSA key used to create digital signatures.
			if skPriv == nil {
				skPriv, err = rsa.GenerateKey(rand.Reader, 2048)
			}
		}
	}
	// The mock client uses a system-assigned endpoint
	if err == nil {
		var endPoint *xt.TcpEndPoint
		endPoint, err = xt.NewTcpEndPoint("127.0.0.1:0")
		if err == nil {
			ep = []xt.EndPointI{endPoint}
		}
	}
	// spin up an XLattice node
	if err == nil {
		node, err = xn.New(name, nil, // get a default NodeID
			lfs, ckPriv, skPriv, nil, ep, nil) // nil overlays, peers
	}
	if err == nil {
		uc, err = NewUpaxClient(ckPriv, skPriv, node, members)
		if err == nil {
			mc = &MockUpaxClient{UpaxClient: *uc}
		}
	}
	return
}

func (muc *MockUpaxClient) createData() (err error) {
	// XXX STUB
	return
}

func (muc *MockUpaxClient) postData() (err error) {
	// XXX STUB
	return
}

func (muc *MockUpaxClient) checkPrimaryServer() (err error) {
	// XXX STUB
	return
}

func (muc *MockUpaxClient) checkOtherServers() (err error) {
	// XXX STUB
	return
}
