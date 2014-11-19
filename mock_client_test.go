package upax_go

import (
	"crypto/rand"
	"crypto/rsa"
	xr "github.com/jddixon/rnglib_go"
	xcl "github.com/jddixon/xlCluster_go"
	xn "github.com/jddixon/xlNode_go"
	xt "github.com/jddixon/xlTransport_go"
	xf "github.com/jddixon/xlUtil_go/lfs"
)

// upax_go/mock_upax_client_test.go

type MockUpaxClient struct {
	K3      int // number of data items
	L1, L2  int // min and max length thereof
	data    [][]byte
	primary uint // which server we will be using
	UpaxClient
}

func NewMockUpaxClient(name, lfs string, members []*xcl.MemberInfo,
	primary uint) (mc *MockUpaxClient, err error) {
	var (
		ckPriv, skPriv *rsa.PrivateKey
		ep             []xt.EndPointI
		node           *xn.Node
		uc             *UpaxClient
	)

	// lfs should be a well-formed POSIX path; if the directory does
	// not exist we should create it.
	err = xf.CheckLFS(lfs, 0750)

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
		uc, err = NewUpaxClient(ckPriv, skPriv, node, members, primary)
		if err == nil {
			mc = &MockUpaxClient{
				UpaxClient: *uc,
			}
		}
	}
	return
}

// Populate the K3 byte slices to be used for testing
func (muc *MockUpaxClient) createData(rng *xr.PRNG, K3, L1, L2 int) (
	err error) {

	muc.K3 = K3
	muc.L1 = L1
	muc.L2 = L2

	muc.data = make([][]byte, K3)
	for i := 0; i < K3; i++ {
		length := L1 + rng.Intn(L2-L1+1) // so L1..L2 inclusive
		muc.data[i] = make([]byte, length)
		rng.NextBytes(muc.data[i])
	}
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
