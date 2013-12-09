package upax_go

// upax_go/client.go

import (
	"crypto/rsa"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	xu "github.com/jddixon/xlattice_go/util"
)

////////////////////////////////////////////////////////
// The model for this is xlattice_go/reg/client_node.go
////////////////////////////////////////////////////////

type UpaxClient struct {
	cnx     xt.ConnectionI
	Members []*reg.MemberInfo
	Version xu.DecimalVersion

	// server side
	Primary   uint           // selects from Members
	serverEnd xt.EndPointI   // convenience
	serverCK  *rsa.PublicKey // convenience

	// this side
	ckPriv, skPriv *rsa.PrivateKey
	ClientCnxHandler
	xn.Node
}

//
func NewUpaxClient(ckPriv, skPriv *rsa.PrivateKey, node *xn.Node,
	members []*reg.MemberInfo, primary uint) (upc *UpaxClient, err error) {

	if ckPriv == nil || skPriv == nil {
		err = NilRSAKey
	} else if node == nil {
		err = NilNode
	} else if members == nil || len(members) == 0 {
		err = NoMembers
	} else if primary >= uint(len(members)) {
		err = PrimaryOutOfRange
	} else {
		upc = &UpaxClient{
			ckPriv:  ckPriv,
			skPriv:  skPriv,
			Node:    *node,
			Members: members,
			Primary: primary,
		}
	}
	return
}

// Enquire as to whether the Upax server has a datum (file) with the
// content key specified.
//
func (upc *UpaxClient) DoYouHave(key *xi.NodeID) (found bool, err error) {

	// XXX STUB XXX

	return
}

// Retrieve from the Upax cluster metadata and the file (datum) with
// the content key specified; returns an error if the datum cannot be
// found or there is a transmission error.
//
func (upc *UpaxClient) Get(key *xi.NodeID) (
	logEntry *LogEntry, data []byte, err error) {

	// XXX STUB XXX

	return
}

// Insert into the Upax cluster metadata and the file described by that
// metadata.  If the data is already present in the cluster the attempt
// to reinsert it will be silently ignored.  If the metatdata is ill-formed
// or does not match the data passed (specifically if it content key
// in the metatdata does not match the data, an error will be returned.
//
func (upc *UpaxClient) Put(logEntry *LogEntry, data []byte) (err error) {

	// XXX STUB XXX

	return
}
