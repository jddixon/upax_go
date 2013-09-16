package node

// upax_go/node/node_test.go

import (
	// "fmt"
	// "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	// "strings"
)

func (s *XLSuite) TestAttrBits(c *C) {

	// client is 00, mirror is 01, server is 10
	c.Assert(UPAX_MIRROR, Equals, 1)
	c.Assert(UPAX_SERVER, Equals, 2)
	c.Assert(WHATEVER, Equals, 4)		// idle frivolity

}
