package upax_go

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/jddixon/xlattice_go/reg"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
)

var _ = fmt.Print

func (s *XLSuite) TestCluster(c *C) {
	rng := xr.MakeSimpleRNG()
	s.doTestCluster(c, rng, true)  // usingSHA1
	s.doTestCluster(c, rng, false) // not
}

func (s *XLSuite) doTestCluster(c *C, rng *xr.PRNG, usingSHA1 bool) {

	// read regCred.dat to get keys etc for a registry --------------
	dat, err := ioutil.ReadFile("regCred.dat")
	c.Assert(err, IsNil)
	regCred, err := reg.ParseRegCred(string(dat))
	c.Assert(err, IsNil)
	regServerName := regCred.Name
	regServerID := regCred.ID
	regServerEnd := regCred.EndPoints[0]
	regServerCK := regCred.CommsPubKey
	regServerSK := regCred.SigPubKey

	// Devise a unique cluster name.  We rely on the convention -----
	// that in Upax tests, the local file system for Upax servers is
	// tmp/CLUSTER-NAME/SERVER-NAME.

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

	// Set the test size in various senses --------------------------
	// K1 is the number of servers, and so the cluster size.  K2 is
	// the number of clients, M the number of messages sent (items to
	// be added to the Upax store), LMin and LMax message lengths.
	K1 := 3 + rng.Intn(5)  // so 3..7
	K2 := 2 + rng.Intn(4)  // so 2..5
	M := 16 + rng.Intn(16) // 16..31
	LMin := 64 + rng.Intn(64)
	LMax := 128 + rng.Intn(128)

	// Use an admin client to get a clusterID for this clusterName --
	const EP_COUNT = 2
	an, err := reg.NewAdminClient(regServerName, regServerID, regServerEnd,
		regServerCK, regServerSK, clusterName, uint64(0), K1, EP_COUNT, nil)
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

	// Create names and LFSs for the K1 servers ---------------------
	// We create a distinct tmp/clusterName/serverName for each
	// server as its local file system (LFS).
	serverNames := make([]string, K1)
	serverPaths := make([]string, K1)
	ckPriv := make([]*rsa.PrivateKey, K1)
	skPriv := make([]*rsa.PrivateKey, K1)
	for i := 0; i < K1; i++ {
		serverNames[i] = rng.NextFileName(8)
		serverPaths[i] = filepath.Join(clusterPath, serverNames[i])
		found, err = xf.PathExists(serverPaths[i])
		c.Assert(err, IsNil)
		for found {
			serverNames[i] = rng.NextFileName(8)
			serverPaths[i] = filepath.Join(clusterPath, serverNames[i])
			found, err = xf.PathExists(serverPaths[i])
			c.Assert(err, IsNil)
		}
		err = os.MkdirAll(serverPaths[i], 0750)
		c.Assert(err, IsNil)
		ckPriv[i], err = rsa.GenerateKey(rand.Reader, 1024) // cheap keys
		c.Assert(err, IsNil)
		c.Assert(ckPriv[i], NotNil)
		skPriv[i], err = rsa.GenerateKey(rand.Reader, 1024) // cheap keys
		c.Assert(err, IsNil)
		c.Assert(skPriv[i], NotNil)
	}

	// create K1 client nodes ---------------------------------------
	uc := make([]*reg.UserClient, K1)
	for i := 0; i < K1; i++ {
		var ep *xt.TcpEndPoint
		ep, err = xt.NewTcpEndPoint("127.0.0.1:0")
		c.Assert(err, IsNil)
		e := []xt.EndPointI{ep}
		uc[i], err = reg.NewUserClient(serverNames[i], serverPaths[i],
			ckPriv[i], skPriv[i],
			regServerName, regServerID, regServerEnd, regServerCK, regServerSK,
			clusterName, cn.ClusterAttrs, cn.ClusterID,
			K1, EP_COUNT, e)
		c.Assert(err, IsNil)
		c.Assert(uc[i], NotNil)
		c.Assert(uc[i].ClusterID, NotNil)
	}
	// Start the K1 client nodes running ----------------------------
	for i := 0; i < K1; i++ {
		err = uc[i].Run()
		c.Assert(err, IsNil)
	}

	// wait until all clientNodes are done --------------------------
	for i := 0; i < K1; i++ {
		<-uc[i].ClientNode.DoneCh
	}

	// convert the client nodes to UpaxServers ----------------------
	us := make([]*UpaxServer, K1)
	for i := 0; i < K1; i++ {
		err = uc[i].PersistClusterMember()
		c.Assert(err, IsNil)
		us[i], err = NewUpaxServer(
			ckPriv[i], skPriv[i], &uc[i].ClusterMember, usingSHA1)
		c.Assert(err, IsNil)
		c.Assert(us[i], NotNil)
	}
	// verify files are present and then start the servers ----------

	// XXX STUB

	// verify servers are running -------------------------

	// XXX STUB

	// When all UpaxServers are ready, create K2 clients.  Each client
	// creates K3 separate datums of differnt length (L1..L2) and
	// content.  Each client signals when done.

	// XXX STUB

	// Verify for each of the K2 clients that its data is present on
	// the selected server.

	// XXX STUB

	// After a reasonable deltaT, verify that all servers have a copy
	// of each and every datum.

	// XXX STUB

	_, _, _, _, _ = K1, K2, M, LMin, LMax
	_ = us
}
