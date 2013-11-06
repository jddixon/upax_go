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
	// member.  XXX Wiring specific values in here is of course bad practice.

	c.Assert(op2tag(XLRegMsg_ItsMe), Equals, MIN_TAG)
	c.Assert(op2tag(XLRegMsg_ItsMe), Equals, 0)

	// KeepAlive msg gets Ack back
	c.Assert(op2tag(XLRegMsg_KeepAlive), Equals, 1)
	c.Assert(op2tag(XLRegMsg_Ack), Equals, 2)

	// Get msg gets Data or Error back
	c.Assert(op2tag(XLRegMsg_Get), Equals, 3)
	c.Assert(op2tag(XLRegMsg_Data), Equals, 4)
	c.Assert(op2tag(XLRegMsg_Error), Equals, 5)

	c.Assert(op2tag(XLRegMsg_IHave), Equals, 6)
	c.Assert(op2tag(XLRegMsg_Put), Equals, 7)

	c.Assert(op2tag(XLRegMsg_Bye), Equals, 8)
	c.Assert(op2tag(XLRegMsg_Bye), Equals, MAX_TAG)
}
