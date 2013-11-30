package upax_go

// upax_go/server.go

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	"github.com/jddixon/xlattice_go/reg"
	xt "github.com/jddixon/xlattice_go/transport"
	"github.com/jddixon/xlattice_go/u"
	xu "github.com/jddixon/xlattice_go/util"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var _ = fmt.Print

const (
	RETRY_INTERVAL = 10 * time.Millisecond
	MAX_RETRY      = 3 // how many times we retry if we can't get a cnx
)

// Load serialized log entries into the IDMap from which they are
// retrievable using the content key.  Conventionally the logEntry
// file is in LFS/U/L, where LFS is the path to the local file system.
//
// We are guaranteed that the file exists and that m is not nil.
//
func loadEntries(pathToFTLog string, m *xn.IDMap, usingSHA1 bool) (
	count int, err error) {

	f, err := os.OpenFile(pathToFTLog, os.O_RDONLY, 0600)
	if err == nil {
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			var entry *LogEntry
			line := scanner.Text()
			entry, err = ParseLogEntry(line, usingSHA1)
			if err != nil {
				break
			}
			err = m.Insert(entry.key, entry)
			if err != nil {
				break
			}
			count++
		}
	}
	return
}

type UpaxServer struct {
	Interval time.Duration
	Lifespan int
	PeerCnx  []xt.ConnectionI
	DoneCh   chan bool

	// conventional log, mostly for debugging
	PathToDebugLog string
	Logger         *log.Logger

	// what we are here for: data managed by the server
	uDir u.UI // content-keyed disk store

	entries     *xn.IDMap // key []byte ==> *LogEntry, stored in U/L
	ftLogFile   *os.File
	pathToFTLog string

	// XXX Should synchronize IDMap internally, not here
	entryCount   int  // number of entries, get lock if changing
	entriesDirty bool // get lock if changing
	entriesMu    sync.RWMutex

	ckPriv, skPriv *rsa.PrivateKey
	reg.ClusterMember
}

func NewUpaxServer(ckPriv, skPriv *rsa.PrivateKey, cm *reg.ClusterMember,
	usingSHA1 bool) (us *UpaxServer, err error) {

	var (
		count     int
		lfs       string   // path to local file system
		f         *os.File // file for debugging log
		pathToLog string
		logger    *log.Logger

		uDir        u.UI
		pathToU     string
		entries     *xn.IDMap
		ftLogFile   *os.File
		pathToFTLog string // conventionally lfs/U/L
	)
	if ckPriv == nil || ckPriv == nil {
		err = NilRSAKey
	} else if cm == nil {
		err = NilClusterMember
	}
	if err == nil {
		serverVersion, err = xu.ParseDecimalVersion(VERSION)
	}
	if err == nil {
		// whatever created cm should have created the local file system
		// and written the node configuration to
		// LFS/.xlattice/cluster.member.config.  Let's make sure that
		// that exists before proceeding.

		lfs = cm.GetLFS()
		// This should be passed in opt.Logger
		pathToLog = filepath.Join(lfs, "log")
		f, err = os.OpenFile(pathToLog,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
		if err == nil {
			logger = log.New(f, "", log.Ldate|log.Ltime)
		}
		pathToCfg := filepath.Join(
			filepath.Join(lfs, ".xlattice"), "cluster.member.config")
		var found bool
		found, err = xf.PathExists(pathToCfg)
		if err == nil && found == false {
			err = ClusterConfigNotFound
		}
	}
	if f != nil {
		defer f.Close()
	}

	if err == nil {
		// DEBUG
		fmt.Printf("creating directory tree in %s\n", lfs)
		// END

		pathToU = filepath.Join(lfs, "U")
		uDir, err = u.New(pathToU, u.DIR16x16, 0)
	}
	if err == nil {
		entries, err = xn.NewNewIDMap() // with default depth
	}
	if err == nil {
		var found bool
		pathToFTLog = filepath.Join(pathToU, "L")
		found, err = xf.PathExists(pathToFTLog)
		if err == nil {
			if found {
				fmt.Printf("ftLog file exists\n")
				count, err = loadEntries(pathToFTLog, entries, usingSHA1)
				if err == nil {
					// reopen it 0600 for appending
					ftLogFile, err = os.OpenFile(pathToFTLog,
						os.O_WRONLY|os.O_APPEND, 0600)
				}
			} else {
				// open it for appending
				ftLogFile, err = os.OpenFile(pathToFTLog,
					os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
			}
		}

	}
	if err == nil {
		us = &UpaxServer{
			DoneCh:         make(chan bool, 2),
			PathToDebugLog: pathToLog,
			Logger:         logger,
			uDir:           uDir,
			entries:        entries,
			ftLogFile:      ftLogFile,
			pathToFTLog:    pathToFTLog,
			entryCount:     count,
			ckPriv:         ckPriv,
			skPriv:         skPriv,
			ClusterMember:  *cm,
		}
	}
	return
}

// Contact a peer, send it an ItsMe, wait for an Ack, and then send
// on the ready chan.  Allow for the fact that the peer may be this
// server, in which case we just send ready.  Also presume the peer
// will take some time to get running, so pause before trying for a
// connection, and be willing to retry.
//
func (us *UpaxServer) InitialHandshake(peerNdx uint32, readyCh chan bool) {

	if peerNdx != us.SelfIndex {
		if peerNdx >= us.ClusterSize {
			panic(fmt.Sprintf("peer index is %d but cluster size is %d",
				peerNdx, us.ClusterSize))
		}
		peerInfo := us.Members[peerNdx]
		peerID := peerInfo.GetNodeID().Value()
		time.Sleep(RETRY_INTERVAL)

		// XXX STUB XXX

		_ = peerID
	}
	readyCh <- true
}

// Run the server, sending keepalives at the interval specified.  The
// server's lifetime is specified in keepalive intervals; if the lifetime
// is less then or equl to zero, it is infinite.
//
func (us *UpaxServer) Run(interval time.Duration, lifespan int) (err error) {

	if interval <= 0 {
		return IntervalMustBePositive
	}
	us.Interval = interval
	us.Lifespan = lifespan

	clusterSize := us.ClusterSize
	us.PeerCnx = make([]xt.ConnectionI, clusterSize)
	serverHasAcked := make([]chan bool, clusterSize)
	for i := uint32(0); i < clusterSize; i++ {
		serverHasAcked[i] = make(chan bool)
	}

	// everything else is in a goroutine
	go func(deltaT time.Duration, lifespan int) {

		for i := uint32(0); i < clusterSize; i++ {

			// Send ItsMe to each other server in the cluster, remembering
			// to skip self.

			// Goroutine signals when Ack is received.

		}

		// XXX STUB XXX

		us.DoneCh <- true
	}(interval, lifespan)

	return
}
