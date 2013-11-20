package upax_go

// upax_go/s_keepalive.go

import (
	"time"
)

type KeepAliveMgr struct {
	Interval time.Duration
	MsgCh    chan *UpaxClusterMsg
	StopCh   chan bool
}

func NewKeepAliveMgr( //h *ClusterOutHandler,
	interval time.Duration, msgCh chan *UpaxClusterMsg,
	stopCh chan bool) (mgr *KeepAliveMgr, err error) {

	if msgCh == nil {
		err = NilMsgCh
	} else {
		mgr = &KeepAliveMgr{
			Interval: interval,
			MsgCh:    msgCh,
			StopCh:   stopCh}
	}
	return
}

func (mgr *KeepAliveMgr) Run() {

	for {
		select {
		case <-time.After(mgr.Interval):
			op := UpaxClusterMsg_KeepAlive
			msgOut := &UpaxClusterMsg{
				Op: &op,
				// MsgN needs to be assigned when the message is
				// actually sent.
			}
			mgr.MsgCh <- msgOut

		case <-mgr.StopCh:
			break
		}
	}
}
