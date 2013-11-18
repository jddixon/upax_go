package upax_go

// upax_go/server.go

import (
	"bufio"
	"crypto/rsa"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	"github.com/jddixon/xlattice_go/reg"
	"github.com/jddixon/xlattice_go/u"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var _ = fmt.Print

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

func (us *UpaxServer) Run() {

	// XXX STUB XXX

}
