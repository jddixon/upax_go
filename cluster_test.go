package upax_go

import (
	"fmt"
	"github.com/jddixon/xlattice_go/reg"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"path/filepath"
)

var _ = fmt.Print

func (s *XLSuite) TestCluster(c *C) {

	// read regCred.dat to get keys etc for a registry
	dat, err := ioutil.ReadFile("regCred.dat")
	c.Assert(err, IsNil)
	regCred, err := reg.ParseRegCred(string(dat))
	c.Assert(err, IsNil)
	serverName := regCred.Name
	serverID := regCred.ID
	serverEnd := regCred.EndPoints[0]
	serverCK := regCred.CommsPubKey
	serverSK := regCred.SigPubKey

	// Devise a unique cluster name.  We rely on the convention that
	// in Upax tests, the local file system for Upax servers is
	// tmp/CLUSTER-NAME/SERVER-NAME.

	rng := xr.MakeSimpleRNG()

	clusterName := rng.NextFileName(8)
	clusterPath := filepath.Join("tmp", clusterName)
	found, err := xf.PathExists(clusterPath)
	c.Assert(err, IsNil)
	for found {
		clusterName = rng.NextFileName(8)
		clusterPath = filepath.Join("tmp", clusterName)
		found, err = xf.PathExists(clusterPath)
		c.Assert(err, IsNil)
	}

	// K1 is the number of servers, and so the cluster size.  K2 is
	// the number of clients, M the number of messages sent (items to
	// be added to the Upax store), LMin and LMax message lengths.
	K1 := 3 + rng.Intn(5)  // so 3..7
	K2 := 2 + rng.Intn(4)  // so 2..5
	M := 16 + rng.Intn(16) // 16..31
	LMin := 64 + rng.Intn(64)
	LMax := 128 + rng.Intn(128)

	// Use an admin client to get a clusterID for this clusterName
	const EP_COUNT = 2
	an, err := reg.NewAdminClient(serverName, serverID, serverEnd,
		serverCK, serverSK, clusterName, uint64(0), K1, EP_COUNT, nil)
	c.Assert(err, IsNil)
	an.Run()
	cn := &an.ClientNode
	<-cn.DoneCh
	clusterID := cn.ClusterID
	c.Assert(clusterID, NotNil)
	clusterSize := cn.ClusterSize
	c.Assert(clusterSize, Equals, uint32(K1))
	epCount := cn.EpCount
	c.Assert(epCount, Equals, uint32(EP_COUNT))

	// Spin up K1 servers in separate goroutines, each with a unique
	// name, each creating tmp/clusterName/serverName as its LFS.

	// XXX STUB

	// When all servers are ready, create K2 clients.  Each client
	// creates K3 separate datums of differnt length (L1..L2) and
	// content.  Each client signals when done.

	// XXX STUB

	// Verify for each of the K2 clients that its data is present on
	// the selected server.

	// XXX STUB

	// After a reasonable deltaT, verify that all servers have a ccopy
	// of each and every datum.

	// XXX STUB

	_, _, _, _, _ = K1, K2, M, LMin, LMax
}
