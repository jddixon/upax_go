package upax_go

// upax_go/c_keepalive.go

import (
	xu "github.com/jddixon/xlattice_go/util"
	"time"
)

// Send keepalives at the specified interval.  If the lifespan is
// greater than 0, do this that many times and then halt.  Otherwise
// do this forever.
//
type ClientKeepAliveMgr struct {
	Interval time.Duration
	lifespan int
	soFar    int
	MsgCh    chan *UpaxClientMsg
	StopCh   chan bool
	DoneCh   chan bool
}

// Send keepalive messages to msgCh at the specified interval.
// The interval should be expressed in units of time, but the lifespan
// is the number of intervals, except that if the lifespan is less
// than or equal to zero, the lifespan is effectively infinite.
//
func NewClientKeepAliveMgr(interval time.Duration, lifespan int,
	msgCh chan *UpaxClientMsg, stopCh, doneCh chan bool) (
	mgr *ClientKeepAliveMgr, err error) {

	if msgCh == nil {
		err = NilMsgCh
	} else {
		if lifespan <= 0 {
			lifespan = xu.MAX_INT
		}
		mgr = &ClientKeepAliveMgr{
			Interval: interval,
			lifespan: lifespan,
			MsgCh:    msgCh,
			StopCh:   stopCh,
			DoneCh:   doneCh}
	}
	return
}

func (mgr *ClientKeepAliveMgr) Run() {

	running := true
	for running {
		select {
		case <-time.After(mgr.Interval):
			op := UpaxClientMsg_KeepAlive
			msgOut := &UpaxClientMsg{
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
