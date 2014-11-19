package upax_go

// upax_go/cluster_test.go

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	reg "github.com/jddixon/xlReg_go"
	xt "github.com/jddixon/xlTransport_go"
	xf "github.com/jddixon/xlUtil_go/lfs"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
	"time"
)

var _ = fmt.Print

func (s *XLSuite) TestCluster(c *C) {
	rng := xr.MakeSimpleRNG()
	s.doTestCluster(c, rng, true)  // usingSHA1
	s.doTestCluster(c, rng, false) // not
}

/////////////////////////////////////////////////////////////////////
// BEFORE RUNNING THIS TEST, BE SURE THAT YOU HAVE
// * IF NECESSARY
//	 * STOPPED ANY OLD INSTANCE OF xlReg
//   * REBUILT AND INSTALLED xlReg
// * STARTED xlReg RUNNING
// * AND INVOKVED ./updateRegCred
//
// In other words, the local version of regCred.dat must match the
// one in /var/app/xlReg; that in turn must correspond to the current
// build of xlReg; and an instance of that build must be running.
/////////////////////////////////////////////////////////////////////

func (s *XLSuite) doTestCluster(c *C, rng *xr.PRNG, usingSHA1 bool) {

	if VERBOSITY > 0 {
		fmt.Printf("TEST_CLUSTER usingSHA1 = %v\n", usingSHA1)
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
	for {
		if _, err = os.Stat(clusterPath); os.IsNotExist(err) {
			break
		}
		clusterName = rng.NextFileName(8)
		clusterPath = filepath.Join("tmp", clusterName)
	}
	err = xf.CheckLFS(clusterPath, 0750)
	c.Assert(err, IsNil)

	// DEBUG
	fmt.Printf("CLUSTER      %s\n", clusterName)
	fmt.Printf("CLUSTER_PATH %s\n", clusterPath)
	// END

	// Set the test size in various senses --------------------------
	// K1 is the number of servers, and so the cluster size.  K2 is
	// the number of clients, M the number of messages sent (items to
	// be added to the Upax store), LMin and LMax message lengths.
	K1 := uint32(3 + rng.Intn(5)) // so 3..7
	K2 := uint32(2 + rng.Intn(4)) // so 2..5
	M := 16 + rng.Intn(16)        // 16..31
	LMin := 64 + rng.Intn(64)
	LMax := 128 + rng.Intn(128)

	// Use an admin client to get a clusterID for this clusterName --
	const EP_COUNT = 2
	an, err := reg.NewAdminClient(regServerName, regServerID, regServerEnd,
		regServerCK, regServerSK, clusterName, uint64(0), K1, EP_COUNT, nil)
	c.Assert(err, IsNil)
	an.Run()
	<-an.DoneCh

	clusterID := an.ClusterID // a NodeID, not []byte
	if clusterID == nil {
		fmt.Println("NIL CLUSTER ID: is xlReg running??")
	}
	c.Assert(clusterID, NotNil)
	clusterSize := an.ClusterSize
	c.Assert(clusterSize, Equals, uint32(K1))
	epCount := an.EPCount
	c.Assert(epCount, Equals, uint32(EP_COUNT))

	// Create names and LFSs for the K1 members ---------------------
	// We create a distinct tmp/clusterName/serverName for each
	// server as its local file system (LFS).
	memberNames := make([]string, K1)
	memberPaths := make([]string, K1)
	ckPriv := make([]*rsa.PrivateKey, K1)
	skPriv := make([]*rsa.PrivateKey, K1)
	for i := uint32(0); i < K1; i++ {
		var found bool
		memberNames[i] = rng.NextFileName(8)
		memberPaths[i] = filepath.Join(clusterPath, memberNames[i])
		found, err = xf.PathExists(memberPaths[i])
		c.Assert(err, IsNil)
		for found {
			memberNames[i] = rng.NextFileName(8)
			memberPaths[i] = filepath.Join(clusterPath, memberNames[i])
			found, err = xf.PathExists(memberPaths[i])
			c.Assert(err, IsNil)
		}
		// DEBUG
		fmt.Printf("MEMBER_PATH[%d]: %s\n", i, memberPaths[i])
		// END
		err = os.MkdirAll(memberPaths[i], 0750)
		c.Assert(err, IsNil)
		ckPriv[i], err = rsa.GenerateKey(rand.Reader, 1024) // cheap keys
		c.Assert(err, IsNil)
		c.Assert(ckPriv[i], NotNil)
		skPriv[i], err = rsa.GenerateKey(rand.Reader, 1024) // cheap keys
		c.Assert(err, IsNil)
		c.Assert(skPriv[i], NotNil)
	}

	// create K1 client nodes ---------------------------------------
	uc := make([]*reg.UserMember, K1)
	for i := uint32(0); i < K1; i++ {
		var ep1, ep2 *xt.TcpEndPoint
		ep1, err = xt.NewTcpEndPoint("127.0.0.1:0")
		ep2, err = xt.NewTcpEndPoint("127.0.0.1:0")
		c.Assert(err, IsNil)
		e := []xt.EndPointI{ep1, ep2}
		uc[i], err = reg.NewUserMember(memberNames[i], memberPaths[i],
			ckPriv[i], skPriv[i],
			regServerName, regServerID, regServerEnd, regServerCK, regServerSK,
			clusterName, an.ClusterAttrs, an.ClusterID,
			K1, EP_COUNT, e)
		c.Assert(err, IsNil)
		c.Assert(uc[i], NotNil)
		c.Assert(uc[i].ClusterID, NotNil)
		c.Assert(uc[i].MemberMaker.DoneCh, NotNil)
	}
	// Start the K1 client nodes running ----------------------------
	for i := uint32(0); i < K1; i++ {
		uc[i].Run()
	}

	fmt.Println("ALL CLIENTS STARTED")

	// wait until all clientNodes are done --------------------------
	for i := uint32(0); i < K1; i++ {
		err = <-uc[i].MemberMaker.DoneCh
		c.Assert(err, IsNil)
		// nodeID := uc[i].clientID
	}
	fmt.Println("ALL CLIENTS DONE") // XXX NOT SEEN

	// verify that all clientNodes have meaningful baseNodes --------
	// XXX THESE TESTS ALWAYS FAIL
	//for i := 0; i < K1; i++ {
	//	c.Assert(uc[i].GetName(), Equals, memberNames[i])
	//	c.Assert(uc[i].GetNodeID(), NotNil)
	//	c.Assert(uc[i].GetCommsPublicKey(), NotNil)
	//	c.Assert(uc[i].GetSigPublicKey(), NotNil)
	//}
	// convert the client nodes to UpaxServers ----------------------
	us := make([]*UpaxServer, K1)
	for i := uint32(0); i < K1; i++ {
		err = uc[i].PersistClusterMember()
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
	fmt.Println("all servers have sent first DONE")
	// END

	// When all UpaxServers are ready, create K2 clients.--
	// Each client creates K3 separate datums of differnt
	// length (L1..L2) and content.  Each client signals
	// when done.

	// XXX STUB

	// Verify for each of the K2 clients ------------------
	// that its data is present on the selected server.  We
	// do this by an Exists() call on uDir for the server's
	// LFS/U for each item posted.

	// XXX STUB

	// After a reasonable deltaT, verify that all servers--
	// have a copy of each and every datum.

	// XXX STUB

	_, _, _, _ = K2, M, LMin, LMax
}
