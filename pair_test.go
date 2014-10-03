package upax_go

// upax_go/pair_test.go

// This is a simplified version of cluster_test.go.  We start a single
// server and then a single client.  The client generates a number of
// dummy files, loads them into the server, and then checks to see
// whether they are there.

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	reg "github.com/jddixon/xlReg_go"
	xr "github.com/jddixon/rnglib_go"
	xt "github.com/jddixon/xlTransport_go"
	xf "github.com/jddixon/xlUtil_go/lfs"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
	"time"
)

var _ = fmt.Print

func (s *XLSuite) TestPair(c *C) {
	rng := xr.MakeSimpleRNG()
	s.doTestPair(c, rng, true)  // usingSHA1
	s.doTestPair(c, rng, false) // not
}

// This was copied from cluster_test.go and minimal changes have been
// made.
//
func (s *XLSuite) doTestPair(c *C, rng *xr.PRNG, usingSHA1 bool) {

	if VERBOSITY > 0 {
		fmt.Printf("TEST_PAIR usingSHA1 = %v\n", usingSHA1)
	}

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
	// K1 is the number of upax servers, and so the cluster size.  K2 is
	// the number of upax clients, M the number of messages sent (items to
	// be added to the Upax store), LMin and LMax message lengths.
	K1 := uint32(2)
	K2 := 1
	M := 16 + rng.Intn(16) // 16..31
	LMin := 64 + rng.Intn(64)
	LMax := 128 + rng.Intn(128)

	// Use an admin client to get a clusterID for this clusterName --
	const EP_COUNT = 2
	an, err := reg.NewAdminMember(regServerName, regServerID, regServerEnd,
		regServerCK, regServerSK, clusterName, uint64(0), K1, EP_COUNT, nil)
	c.Assert(err, IsNil)
	an.Run()
	cn := &an.MemberNode
	<-cn.DoneCh
	clusterID := cn.ClusterID
	if clusterID == nil {
		fmt.Println("NIL CLUSTER ID: is xlReg running??")
	}
	c.Assert(clusterID, NotNil)
	clusterSize := cn.ClusterSize
	c.Assert(clusterSize, Equals, uint32(K1))
	epCount := cn.EpCount
	c.Assert(epCount, Equals, uint32(EP_COUNT))

	// DEBUG
	// fmt.Printf("cluster %s: %s\n", clusterName, clusterID.String())
	// END

	// Create names and LFSs for the K1 servers ---------------------
	// We create a distinct tmp/clusterName/serverName for each
	// server as its local file system (LFS).
	serverNames := make([]string, K1)
	serverPaths := make([]string, K1)
	ckPriv := make([]*rsa.PrivateKey, K1)
	skPriv := make([]*rsa.PrivateKey, K1)
	for i := uint32(0); i < K1; i++ {
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

	// create K1 reg client nodes -----------------------------------
	uc := make([]*reg.UserMember, K1)
	for i := uint32(0); i < K1; i++ {
		var ep *xt.TcpEndPoint
		ep, err = xt.NewTcpEndPoint("127.0.0.1:0")
		c.Assert(err, IsNil)
		e := []xt.EndPointI{ep}
		uc[i], err = reg.NewUserMember(serverNames[i], serverPaths[i],
			ckPriv[i], skPriv[i],
			regServerName, regServerID, regServerEnd, regServerCK, regServerSK,
			clusterName, cn.ClusterAttrs, cn.ClusterID,
			K1, EP_COUNT, e)
		c.Assert(err, IsNil)
		c.Assert(uc[i], NotNil)
		c.Assert(uc[i].ClusterID, NotNil)
	}
	// Start the K1 reg client nodes running ------------------------
	for i := uint32(0); i < K1; i++ {
		uc[i].Run()
	}

	// wait until all reg clientNodes are done ----------------------
	for i := uint32(0); i < K1; i++ {
		err := <-uc[i].MemberNode.DoneCh
		c.Assert(err, IsNil)
	}

	// verify that all clientNodes have meaningful baseNodes --------
	//for i := 0; i < K1; i++ {
	//	c.Assert(uc[i].GetName(), Equals, serverNames[i])
	//	c.Assert(uc[i].GetNodeID(), NotNil)
	//	c.Assert(uc[i].GetCommsPublicKey(), NotNil)
	//	c.Assert(uc[i].GetSigPublicKey(), NotNil)
	//}

	// verify that all clientNode members have meaningful baseNodes -
	for i := uint32(0); i < K1; i++ {
		// fmt.Printf("  server %s\n", serverNames[i])	// DEBUG
		memberCount := len(uc[i].Members)
		c.Assert(memberCount, Equals, K1)
		for j := 0; j < memberCount; j++ {
			c.Assert(uc[i].Members[j], NotNil)
			// DEBUG
			// fmt.Printf("    other server[%d] is %s\n", j, serverNames[j])
			// END

			// doesn't work because reg server does not necessarily see
			// members in serverName order.
			// c.Assert(uc[i].Members[j].GetName(), Equals, serverNames[j])
			c.Assert(uc[i].Members[j].GetName() == "", Equals, false)
			c.Assert(uc[i].Members[j].GetNodeID(), NotNil)
			c.Assert(uc[i].Members[j].GetCommsPublicKey(), NotNil)
			c.Assert(uc[i].Members[j].GetSigPublicKey(), NotNil)
		}
	}

	// convert the reg client nodes to UpaxServers ------------------
	us := make([]*UpaxServer, K1)
	for i := uint32(0); i < K1; i++ {
		err = uc[i].PersistClusterMember() // sometimes panics
		c.Assert(err, IsNil)
		us[i], err = NewUpaxServer(
			ckPriv[i], skPriv[i], &uc[i].ClusterMember, usingSHA1)
		c.Assert(err, IsNil)
		c.Assert(us[i], NotNil)
	}
	// verify files are present and then start the servers ----------

	// 11-07 TODO, modified:
	// Run() causes each server to send ItsMe to all other servers;
	// as each gets its Ack, it starts the KeepAlive/Ack cycle running
	// at a 50 ms interval specified as a Run() argument and then sends
	// on DoneCh.  Second parameter is lifetime of the server in
	// keep-alives, say 20 (so 1 sec in total).  When this time has
	// passed, the server will send again on DoneCh, and then shut down.

	// XXX STUB
	for i := uint32(0); i < K1; i++ {
		err = us[i].Run(10*time.Millisecond, 20)
		c.Assert(err, IsNil)
	}

	// Verify servers are running -------------------------
	// 11-18: we wait for the first done from each server.
	//
	// XXX STUB
	for i := uint32(0); i < K1; i++ {
		<-us[i].DoneCh
	}
	// DEBUG
	fmt.Println("pair_test: both servers have sent first DONE")
	// END

	// When all UpaxServers are ready, create K2 clients.--
	// Each upax client creates K3 separate datums of different
	// length (L1..L2) and content.  Each client signals
	// when done.

	// XXX STUB

	// Verify for each of the K2 clients ------------------
	// that its data is present on the selected server.  We
	// do this by an Exists() call on uDir for the server's
	// LFS/U for each item posted.

	// XXX STUB

	// After a reasonable deltaT, verify that both servers--
	// have a copy of each and every datum.

	// XXX STUB

	_, _, _, _ = K2, M, LMin, LMax
}
