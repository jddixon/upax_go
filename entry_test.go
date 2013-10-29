package upax_go

// upax_go/entry_test.go

import (
	"encoding/hex"
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"strings"
)

func (s *XLSuite) doTestEntry(c *C, rng *rnglib.PRNG, usingSHA1 bool) {
	t := rng.Int63()
	var length int
	if usingSHA1 {
		length = 20
	} else {
		length = 32
	}
	key := make([]byte, length)
	rng.NextBytes(&key)
	hexKey := hex.EncodeToString(key)

	nodeID := make([]byte, length)
	rng.NextBytes(&nodeID)
	hexNodeID := hex.EncodeToString(nodeID)

	src := rng.NextFileName(32) // 32 is max len
	path := rng.NextFileName(32)
	for strings.Contains(path, ".") { // that crude fix
		path = rng.NextFileName(32)
	}

	expected := fmt.Sprintf("%d %s %s \"%s\" %s",
		t, hexKey, hexNodeID, src, path)

	entry, err := NewLogEntry(t, key, nodeID, src, path)
	c.Assert(err, IsNil)

	c.Assert(entry.TimeStamp(), Equals, t)
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
