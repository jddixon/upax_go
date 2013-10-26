package node

// upax_go/node/node_test.go

import (
	// "fmt"
	// "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	// "strings"
)

func (s *XLSuite) TestAttrBits(c *C) {

	// A silly test, and wasteful: we use three bits to distinguish
	// three states.
	// 
	// Client is 01, mirror is 10, server is 100.
	c.Assert(UPAX_CLIENT, Equals, RoleBits(1))
	c.Assert(UPAX_MIRROR, Equals, RoleBits(2))
	c.Assert(UPAX_SERVER, Equals, RoleBits(4))

}
