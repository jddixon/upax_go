package upax_go

// upax_go/s_ihave_mgr.go

import (
	xi "github.com/jddixon/xlNodeID_go"
)

// The ClusterIHaveMgr receives IHaveObjs on its input channel.  For each
// ID, it checks to see whether the ID is present in its entries.
// If the ID is not present, it creates a Get message, which it passes
// on.
//
type ClusterIHaveMgr struct {
	iHaveCh  chan IHaveObj
	entries  *xi.IDMap // a convenience
	outMsgCh chan *UpaxClusterMsg
	stopCh   chan bool
}

func NewClusterIHaveMgr(iHaveCh chan IHaveObj, entries *xi.IDMap,
	outMsgCh chan *UpaxClusterMsg, stopCh chan bool) (
	mgr *ClusterIHaveMgr, err error) {

	if iHaveCh == nil {
		err = NilIHaveChan
	} else if entries == nil {
		err = NilIDMap
	} else if outMsgCh == nil {
		err = NilOutMsgCh
	} else {
		mgr = &ClusterIHaveMgr{
			iHaveCh:  iHaveCh,
			entries:  entries,
			outMsgCh: outMsgCh,
			stopCh:   stopCh,
		}
	}
	return
}

// This will normally be run in a separate goroutine.
//
func (mgr *ClusterIHaveMgr) Run() {
	var whatever interface{}
	var err error
	for {
		select {
		case iHaveObj := <-mgr.iHaveCh:
			ids := iHaveObj.IDs
			for i := 0; i < len(ids); i++ {
				id := ids[i]
				whatever, err = mgr.entries.Find(id)
				if (err == nil) && (whatever == nil) {
					op := UpaxClusterMsg_Get
					msgOut := &UpaxClusterMsg{
						Op:   &op,
						Hash: id,
						// MsgN needs to be assigned when the message is
						// actually sent.
					}
					mgr.outMsgCh <- msgOut
				}
			}
		case <-mgr.stopCh:
			break
		}
	}
}
