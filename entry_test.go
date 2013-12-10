package upax_go

// upax_go/entry_test.go

import (
	"bytes"
	"encoding/hex"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	"github.com/jddixon/xlattice_go/rnglib"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	. "launchpad.net/gocheck"
	"os"
	"path/filepath"
	"strings"
)

// for loadEntries()
import ()

func (s *XLSuite) makeEntryData(c *C, rng *rnglib.PRNG, usingSHA1 bool) (
	t int64, key, nodeID []byte, src, path string) {

	t = rng.Int63() // timestamp
	var length int
	if usingSHA1 {
		length = 20
	} else {
		length = 32
	}
	key = make([]byte, length)
	rng.NextBytes(&key)

	nodeID = make([]byte, length)
	rng.NextBytes(&nodeID)

	src = rng.NextFileName(32) // 32 is max len
	path = rng.NextFileName(32)
	for strings.Contains(path, ".") { // that crude fix
		path = rng.NextFileName(32)
	}
	return
}
func (s *XLSuite) doTestEntry(c *C, rng *rnglib.PRNG, usingSHA1 bool) {

	t, key, nodeID, src, path := s.makeEntryData(c, rng, usingSHA1)
	hexKey := hex.EncodeToString(key)
	hexNodeID := hex.EncodeToString(nodeID)

	expected := fmt.Sprintf("%d %s %s \"%s\" %s",
		t, hexKey, hexNodeID, src, path)

	entry, err := NewLogEntry(t, key, nodeID, src, path)
	c.Assert(err, IsNil)

	c.Assert(entry.Timestamp(), Equals, t)
	c.Assert(hex.EncodeToString(entry.Key()), Equals, hexKey)
	c.Assert(hex.EncodeToString(entry.NodeID()), Equals, hexNodeID)
	c.Assert(entry.Src(), Equals, src)
	c.Assert(entry.Path(), Equals, path)

	serialization := entry.String()
	c.Assert(serialization, Equals, expected)

	backAgain, err := ParseLogEntry(serialization, usingSHA1)
	c.Assert(err, IsNil)
	reserialization := backAgain.String()
	c.Assert(reserialization, Equals, serialization)
}

func (s *XLSuite) TestEntry(c *C) {
	rng := rnglib.MakeSimpleRNG()
	s.doTestEntry(c, rng, true)
	s.doTestEntry(c, rng, false)
}

// Test the function used by the server to load log entries from
// the disk.  These are conventionally stored in lsf/U/L, but here
// are stored in a randomly named file under ./tmp/
//
func (s *XLSuite) TestLoadEntries(c *C) {
	rng := rnglib.MakeSimpleRNG()
	s.doTestLoadEntries(c, rng, true)  // using SHA1
	s.doTestLoadEntries(c, rng, false) // not using SHA1
}
func (s *XLSuite) doTestLoadEntries(c *C, rng *rnglib.PRNG, usingSHA1 bool) {
	K := 16 + rng.Intn(16)

	// create a unique name for a scratch file
	pathToFile := filepath.Join("tmp", rng.NextFileName(16))
	found, err := xf.PathExists(pathToFile)
	c.Assert(err, IsNil)
	for found {
		pathToFile = filepath.Join("tmp", rng.NextFileName(16))
		found, err = xf.PathExists(pathToFile)
		c.Assert(err, IsNil)
	}
	f, err := os.OpenFile(pathToFile, os.O_CREATE|os.O_WRONLY, 0600)
	c.Assert(err, IsNil)

	// create K entries, saving them in a slice while writing them
	// to disk
	var entries []*LogEntry
	for i := 0; i < K; i++ {
		t, key, nodeID, src, path := s.makeEntryData(c, rng, usingSHA1)
		entry, err := NewLogEntry(t, key, nodeID, src, path)
		c.Assert(err, IsNil)
		strEntry := entry.String()
		entries = append(entries, entry)
		var count int
		count, err = f.WriteString(strEntry + "\n")
		c.Assert(err, IsNil)
		c.Assert(count, Equals, len(strEntry)+1)
	}
	f.Close()
	c.Assert(len(entries), Equals, K)

	// use UpaxServer.LoadEntries to load the stuff in the file.
	m, err := xn.NewNewIDMap()
	c.Assert(err, IsNil)
	count, err := loadEntries(pathToFile, m, usingSHA1)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, K) // K entries loaded.

	for i := 0; i < K; i++ {
		var entry, eInMap *LogEntry
		var whatever interface{}
		entry = entries[i]
		key := entry.key
		whatever, err = m.Find(key)
		c.Assert(err, IsNil)
		c.Assert(whatever, NotNil)
		eInMap = whatever.(*LogEntry)

		// DEBUG
		// XXX NEED LogEntry.Equal()
		// END

		c.Assert(bytes.Equal(key, eInMap.key), Equals, true)
	}
}
