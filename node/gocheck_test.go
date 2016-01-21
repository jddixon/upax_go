package node

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xcl "github.com/jddixon/xlCluster_go"
	xi "github.com/jddixon/xlNodeID_go"
	xn "github.com/jddixon/xlNode_go"
	reg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	xu "github.com/jddixon/xlUtil_go"
	. "gopkg.in/check.v1"
	"strings"
	"testing"
)

// IF USING gocheck, need a file with the next 3 lines in each directory.
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

const (
	VERBOSITY = 1
)

// UTILITY FUNCTIONS FOR THIS PACKAGE ///////////////////////////////

// XXX With the reorganization of the package and changed scheme for
// providing unique nodeIDs (xlReg), much of whese functions are not
// going to be used here.

func (s *XLSuite) makeAnID(c *C, rng *xr.PRNG) (id []byte) {
	id = make([]byte, xu.SHA3_BIN_LEN)
	rng.NextBytes(id)
	return
}
func (s *XLSuite) makeANodeID(c *C, rng *xr.PRNG) (nodeID *xi.NodeID) {
	id := s.makeAnID(c, rng)
	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))
	return
}
func (s *XLSuite) makeAnRSAKey(c *C) (key *rsa.PrivateKey) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)
	c.Assert(key, Not(IsNil))
	return key
}

// Creates a local (127.0.0.1) endPoint and adds it to the node.
// XXX This code was hacked from ../node/node_test.go.

func (s *XLSuite) makeALocalEndPoint(c *C, node *xn.Node) {
	addr := fmt.Sprintf("127.0.0.1:0")
	ep, err := xt.NewTcpEndPoint(addr)
	c.Assert(err, IsNil)
	c.Assert(ep, Not(IsNil))
	ndx, err := node.AddEndPoint(ep)
	c.Assert(err, IsNil)
	c.Assert(ndx, Equals, 0) // it's the only one
}

// Return an initialized and tested host, with a NodeID, commsKey,
// and sigKey.   XXX This code was hacked from ../node/node_test.go
// and then simplified a bit.

func (s *XLSuite) makeHostAndKeys(c *C, rng *xr.PRNG) (
	n *xn.Node, ckPriv, skPriv *rsa.PrivateKey) {

	// XXX names may not be unique
	name := rng.NextFileName(6)
	for {
		first := string(name[0])
		if !strings.Contains(first, "0123456789") &&
			!strings.Contains(name, "-") {
			break
		}
		name = rng.NextFileName(6)
	}
	id := s.makeANodeID(c, rng)
	lfs := "tmp/" + hex.EncodeToString(id.Value())
	ckPriv = s.makeAnRSAKey(c)
	skPriv = s.makeAnRSAKey(c)

	n, err2 := xn.New(name, id, lfs, ckPriv, skPriv, nil, nil, nil)

	c.Assert(err2, IsNil)
	c.Assert(n, Not(IsNil))
	c.Assert(name, Equals, n.GetName())
	actualID := n.GetNodeID()
	c.Assert(true, Equals, id.Equal(actualID))
	// s.doKeyTests(c, n, rng)
	c.Assert(0, Equals, (*n).SizePeers())
	c.Assert(0, Equals, (*n).SizeOverlays())
	c.Assert(0, Equals, n.SizeConnections())
	c.Assert(lfs, Equals, n.GetLFS())
	return n, ckPriv, skPriv
}

// Using functions must check to ensure members have unique names

func (s *XLSuite) makeAMemberInfo(c *C, rng *xr.PRNG) *xcl.MemberInfo {
	attrs := uint64(rng.Int63())
	bn, err := xn.NewBaseNode(
		rng.NextFileName(8),
		s.makeANodeID(c, rng),
		&s.makeAnRSAKey(c).PublicKey,
		&s.makeAnRSAKey(c).PublicKey,
		nil) // overlays
	c.Assert(err, IsNil)

	asPeer := &xn.Peer{
		BaseNode: *bn,
	}
	return &xcl.MemberInfo{
		Attrs: attrs,
		Peer:  asPeer,
	}
}

// Make a RegCluster for test purposes.  Cluster member names are guaranteed
// to be unique but the name of the cluster itself may not be.

func (s *XLSuite) makeACluster(c *C, rng *xr.PRNG, epCount, size uint32) (
	rc *reg.RegCluster) {

	var err error
	c.Assert(1 < size && size <= 64, Equals, true)

	attrs := uint64(rng.Int63())
	name := rng.NextFileName(8) // no guarantee of uniqueness
	id := s.makeANodeID(c, rng)

	rc, err = reg.NewRegCluster(name, id, attrs, size, epCount)
	c.Assert(err, IsNil)

	for count := uint32(0); count < size; count++ {
		cm := s.makeAMemberInfo(c, rng)
		for {
			if _, ok := rc.MembersByName[cm.Peer.GetName()]; ok {
				// name is in use, so try again
				cm = s.makeAMemberInfo(c, rng)
			} else {
				// copy the connector list as strings
				myEnds := make([]string, 0, 0)
				for endCount := 0; endCount < int(size); endCount++ {
					thisEnd := cm.Peer.GetConnector(endCount).String()
					myEnds = append(myEnds, thisEnd)
				}
				asClient := &reg.ClientInfo{
					Attrs:    cm.Attrs,
					MyEnds:   myEnds,
					BaseNode: cm.Peer.BaseNode,
				}
				err = rc.AddMember(asClient)
				break
			}
		}
	}
	return
}
