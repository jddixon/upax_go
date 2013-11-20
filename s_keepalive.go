package upax_go

// upax_go/s_keepalive.go

import (
	"time"
)

type KeepAliveMgr struct {
	Interval time.Duration
	OutMsgCh chan *UpaxClusterMsg
	StopCh   chan bool
}

func NewKeepAliveMgr(h *ClusterOutHandler,
	interval time.Duration, outMsgCh chan *UpaxClusterMsg,
	stopCh chan bool) (mgr *KeepAliveMgr, err error) {

	if outMsgCh == nil {
		err = NilOutMsgCh
	} else {
		mgr = &KeepAliveMgr{
			Interval: interval,
			OutMsgCh: outMsgCh,
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
			mgr.OutMsgCh <- msgOut

		case <-mgr.StopCh:
			break
		}
	}
}
