package upax_go

// upax_go/c_ihave_mgr_test.go

import (
	"fmt"
	xr "github.com/jddixon/rnglib_go"
	xi "github.com/jddixon/xlNodeID_go"
	. "gopkg.in/check.v1"
	"time"
)

func (s *XLSuite) TestClientClientIHaveMgr(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLIENT_IHAVE_MGR")
	}
	rng := xr.MakeSimpleRNG()
	iHaveCh := make(chan IHaveObj)
	entries, err := xi.NewNewIDMap()
	c.Assert(err, IsNil)
	outMsgCh := make(chan *UpaxClientMsg, 16)
	stopCh := make(chan bool)

	K := 3 + rng.Intn(14)
	keys := make([][]byte, K)
	for i := 0; i < K; i++ {
		keys[i] = make([]byte, 32)
		rng.NextBytes(keys[i])
		if i < K/2 {
			err = entries.Insert(keys[i], &keys[i])
			c.Assert(err, IsNil)
		}
	}
	obj := IHaveObj{keys}

	mgr, err := NewClientIHaveMgr(iHaveCh, entries, outMsgCh, stopCh)
	c.Assert(err, IsNil)
	go mgr.Run()
	mgr.iHaveCh <- obj

	var msgs []*UpaxClientMsg

	done := false
	for !done {
		select {
		case msg := <-outMsgCh:
			msgs = append(msgs, msg)
		case <-time.After(time.Millisecond):
			done = true
			break
		}
	}
	c.Assert(len(msgs), Equals, K-K/2)
	stopCh <- true
}
