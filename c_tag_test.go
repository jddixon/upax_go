package upax_go

// upax_go/c_tag_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestClientInHandler(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLIENT_IN_HANDLER")
	}

	// These are the tags that InHandler will accept from another cluster
	// member.  XXX Wiring specific values in here is of course bad practice.

	c.Assert(clientOp2tag(UpaxClientMsg_Intro), Equals, C_MIN_TAG)
	c.Assert(clientOp2tag(UpaxClientMsg_Intro), Equals, uint(0))
	c.Assert(clientOp2tag(UpaxClientMsg_ItsMe), Equals, uint(1))
	c.Assert(clientOp2tag(UpaxClientMsg_KeepAlive), Equals, uint(2))
	c.Assert(clientOp2tag(UpaxClientMsg_Query), Equals, uint(3))
	c.Assert(clientOp2tag(UpaxClientMsg_Get), Equals, uint(4))
	c.Assert(clientOp2tag(UpaxClientMsg_IHave), Equals, uint(5))
	c.Assert(clientOp2tag(UpaxClientMsg_Put), Equals, uint(6))
	c.Assert(clientOp2tag(UpaxClientMsg_Bye), Equals, uint(7))
	c.Assert(clientOp2tag(UpaxClientMsg_Bye), Equals, C_MAX_TAG)

	// These are reply tags.

	c.Assert(clientOp2tag(UpaxClientMsg_Ack), Equals, uint(10))
	c.Assert(clientOp2tag(UpaxClientMsg_Data), Equals, uint(11))
	c.Assert(clientOp2tag(UpaxClientMsg_NotFound), Equals, uint(12))
	c.Assert(clientOp2tag(UpaxClientMsg_Error), Equals, uint(13))

}
