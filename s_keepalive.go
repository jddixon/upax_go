package upax_go

// upax_go/s_keepalive.go

import (
	xu "github.com/jddixon/xlUtil_go"
	"time"
)

// Send keepalives at the specified interval.  If the lifespan is
// greater than 0, do this that many times and then halt.  Otherwise
// do this forever.
//
type ClusterKeepAliveMgr struct {
	Interval time.Duration
	lifespan int
	soFar    int
	MsgCh    chan *UpaxClusterMsg
	StopCh   chan bool
	DoneCh   chan bool
}

// Send keepalive messages to msgCh at the specified interval.
// The interval should be expressed in units of time, but the lifespan
// is the number of intervals, except that if the lifespan is less
// than or equal to zero, the lifespan is effectively infinite.
//
func NewClusterKeepAliveMgr(interval time.Duration, lifespan int,
	msgCh chan *UpaxClusterMsg, stopCh, doneCh chan bool) (
	mgr *ClusterKeepAliveMgr, err error) {

	if msgCh == nil {
		err = NilMsgCh
	} else {
		if lifespan <= 0 {
			lifespan = xu.MAX_INT
		}
		mgr = &ClusterKeepAliveMgr{
			Interval: interval,
			lifespan: lifespan,
			MsgCh:    msgCh,
			StopCh:   stopCh,
			DoneCh:   doneCh}
	}
	return
}

func (mgr *ClusterKeepAliveMgr) Run() {

	running := true
	for running {
		select {
		case <-time.After(mgr.Interval):
			op := UpaxClusterMsg_KeepAlive
			msgOut := &UpaxClusterMsg{
				Op: &op,
				// MsgN needs to be assigned when the message is
				// actually sent.
			}
			mgr.MsgCh <- msgOut
			mgr.soFar++
			if mgr.soFar >= mgr.lifespan {
				running = false
			}

		case <-mgr.StopCh:
			running = false
		}
	}
	// without a delay we may lose the last keepalive)
	time.Sleep(mgr.Interval)
	mgr.DoneCh <- true
}
