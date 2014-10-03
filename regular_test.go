package upax_go

// upax_go/regular_test.go

import (
	"encoding/hex"
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "launchpad.net/gocheck"
	"strings"
)

// ORIGINAL OF THIS FILE WAS OVERWRITTEN - REPAIRS IN PROGRESS

func (s *XLSuite) TestBasics(c *C) {
	c.Assert(IgnorableRE(), Equals, ignorableRE)
	c.Assert(PathRE(), Equals, pathRE)
	c.Assert(BodyLine1RE(), Equals, bodyLine1RE)
	c.Assert(BodyLine3RE(), Equals, bodyLine3RE)
}
func (s *XLSuite) TestIgnorability(c *C) {
	c.Assert(ignorableRE.MatchString("    "), Equals, true)
	c.Assert(ignorableRE.MatchString("  # 123"), Equals, true)
	c.Assert(ignorableRE.MatchString("  // 123"), Equals, false)
}

func (s *XLSuite) TestPathRE(c *C) {
	// XXX STUB XXX
}
func (s *XLSuite) doTestRegexes(c *C, rng *xr.PRNG, usingSHA1 bool) {
	t := rng.Int63()
	var length int
	if usingSHA1 {
		length = 20
	} else {
		length = 32
	}
	key := make([]byte, length)
	rng.NextBytes(key)
	hexKey := hex.EncodeToString(key)

	nodeID := make([]byte, length)
	rng.NextBytes(nodeID)
	hexNodeID := hex.EncodeToString(nodeID)

	src := rng.NextFileName(32) // 32 is max len
	path := rng.NextFileName(32)
	for strings.Contains(path, ".") { // XXX a crude fix
		path = rng.NextFileName(32)
	}

	expected := fmt.Sprintf("%d %s %s \"%s\" %s",
		t, hexKey, hexNodeID, src, path)

	if usingSHA1 {
		c.Assert(bodyLine1RE.MatchString(expected), Equals, true)
		groups := bodyLine1RE.FindStringSubmatch(expected)
		c.Assert(groups, Not(IsNil))
		c.Assert(len(groups), Equals, 6) // 5 fields + match on all

		c.Assert(bodyLine3RE.MatchString(expected), Equals, false)

	} else {
		// DEBUG
		if !bodyLine3RE.MatchString(expected) {
			fmt.Printf("DOESN'T MATCH PATTERN: %s\n", expected)
		}
		// END
		c.Assert(bodyLine3RE.MatchString(expected), Equals, true)
		groups := bodyLine3RE.FindStringSubmatch(expected)
		c.Assert(groups, Not(IsNil))
		c.Assert(len(groups), Equals, 6)

		c.Assert(bodyLine1RE.MatchString(expected), Equals, false)
	}
}

func (s *XLSuite) TestRegexes(c *C) {
	rng := xr.MakeSimpleRNG()
	for i := 0; i < 8; i++ {
		s.doTestRegexes(c, rng, true)
		s.doTestRegexes(c, rng, false)
	}
}
