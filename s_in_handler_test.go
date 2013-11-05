package upax_go

// upax_go/s_in_handler_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestInHandler(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_IN_HANDLER")
	}

	// These are the tags that InHandler will accept from another cluster
	// member.

	c.Assert(op2tag(XLRegMsg_KeepAlive), Equals, MIN_TAG)

	c.Assert(op2tag(XLRegMsg_KeepAlive), Equals, 3)
	c.Assert(op2tag(XLRegMsg_Ack), Equals, 4)
	
	c.Assert(op2tag(XLRegMsg_Error), Equals, 5)

	c.Assert(op2tag(XLRegMsg_Get), Equals, 6)
	c.Assert(op2tag(XLRegMsg_Data), Equals, 7)
	
	
	c.Assert(op2tag(XLRegMsg_Bye), Equals, 13)
	c.Assert(op2tag(XLRegMsg_Bye), Equals, MAX_TAG)
}
