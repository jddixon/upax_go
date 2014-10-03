package upax_go

// upax_go/s_keepalive_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	. "launchpad.net/gocheck"
	"time"
)

func (s *XLSuite) TestClusterKeepAlive(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_KEEP_ALIVE")
	}
	rng := xr.MakeSimpleRNG()
	k := time.Duration(1 + rng.Intn(10))
	interval := k * time.Millisecond
	lifespan := 3 + rng.Intn(13)

	msgCh := make(chan *UpaxClusterMsg, 2*lifespan)
	stopCh := make(chan bool)
	doneCh := make(chan bool)
	var msgs []*UpaxClusterMsg

	mgr, err := NewClusterKeepAliveMgr(
		interval, lifespan, msgCh, stopCh, doneCh)
	c.Assert(err, IsNil)
	go mgr.Run()

	done := false

	select {
	case <-time.After(time.Duration(2*lifespan) * interval):
		c.Fatal("timed out waiting for done from ClusterKeepAliveMgr")
	default:
		for !done {
			select {
			case msg := <-msgCh:
				msgs = append(msgs, msg)
			case <-doneCh:
				done = true
			}
		}
	}
	if done {
		c.Assert(len(msgs), Equals, lifespan)
	} else {
		stopCh <- true
		<-doneCh
	}
}
