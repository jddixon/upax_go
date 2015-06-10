package node

import (
	"crypto/rsa"

	xn "github.com/jddixon/xlNode_go"
	xi "github.com/jddixon/xlNodeID_go"
	xo "github.com/jddixon/xlOverlay_go"
	xt "github.com/jddixon/xlTransport_go"

	"strings"
)

type RoleBits uint64

const (
	UPAX_CLIENT RoleBits = 1 << iota
	UPAX_MIRROR
	UPAX_SERVER
)

type UpaxNode struct {
	Attrs  uint64 // must specify server or mirror
	Acc    xt.AcceptorI
	ckPriv *rsa.PrivateKey
	skPriv *rsa.PrivateKey
	xn.Node
}

// Constructor must specify LFS, which must have restricted access,
// because private keys will be stored there.

type NodeOptions struct {
	Lfs string
}

func New(name string, id *xi.NodeID, lfs string,
	commsKey, sigKey *rsa.PrivateKey,
	o []xo.OverlayI, e []xt.EndPointI, p []*xn.Peer) (
	uNode *UpaxNode, err error) {

	n, err := xn.New(name, id, lfs, commsKey, sigKey, o, e, p)

	if err == nil {
		uNode = &UpaxNode{
			ckPriv: commsKey,
			skPriv: sigKey,
			Node:   *n,
		}
	}
	return
}

// SERIALIZATION ////////////////////////////////////////////////////

func (un *UpaxNode) Strings() []string {

	// XXX STUB XXX

	// first serialize the node, then add the bits specific to UpaxNodes

	return nil
}

func (un *UpaxNode) String() string {
	return strings.Join(un.Strings(), "\n")
}

func Parse(s string) (uNode *UpaxNode, rest []string, err error) {

	// XXX STUB XXX

	// First parse the serialized node, then the bits specific to
	// upaxNodes, leaving anything that can't be handled in rest.

	// Current versions of similar code create a live node, which
	// has an open acceptor - which can cause problems.  Think it
	// through and modify behavior if appropriate.

	return
}
