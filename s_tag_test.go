package upax_go

// upax_go/s_tag_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestClusterInHandler(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_IN_HANDLER")
	}

	// These are the tags that InHandler will accept from another cluster
	// member.  XXX Wiring specific values in here is of course bad practice.

	c.Assert(clusterOp2tag(UpaxClusterMsg_ItsMe), Equals, S_MIN_TAG)
	c.Assert(clusterOp2tag(UpaxClusterMsg_ItsMe), Equals, uint(0))
	c.Assert(clusterOp2tag(UpaxClusterMsg_KeepAlive), Equals, uint(1))
	c.Assert(clusterOp2tag(UpaxClusterMsg_Get), Equals, uint(2))
	c.Assert(clusterOp2tag(UpaxClusterMsg_IHave), Equals, uint(3))
	c.Assert(clusterOp2tag(UpaxClusterMsg_Put), Equals, uint(4))
	c.Assert(clusterOp2tag(UpaxClusterMsg_Bye), Equals, uint(5))
	c.Assert(clusterOp2tag(UpaxClusterMsg_Bye), Equals, S_MAX_TAG)

	// These are reply tags.

	c.Assert(clusterOp2tag(UpaxClusterMsg_Ack), Equals, uint(10))
	c.Assert(clusterOp2tag(UpaxClusterMsg_Data), Equals, uint(11))
	c.Assert(clusterOp2tag(UpaxClusterMsg_NotFound), Equals, uint(12))
	c.Assert(clusterOp2tag(UpaxClusterMsg_Error), Equals, uint(13))

}
