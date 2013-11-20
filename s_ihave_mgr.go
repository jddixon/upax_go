package upax_go

import (
	xn "github.com/jddixon/xlattice_go/node"
)

// Each ID is a content key for an entry that the peer claims to have
// in its store.
//
type IHaveObj struct {
	IDs [][]byte
}

// The IHaveMgr receives IHaveObjs on its input channel.  For each
// ID, it checks to see whether the ID is present in its entries.
// If the ID is not present, it creates a Get message, which it passes
// on.
//
type IHaveMgr struct {
	iHaveCh  chan IHaveObj
	entries  *xn.IDMap // a convenience
	outMsgCh chan *UpaxClusterMsg
	stopCh   chan bool
}

func NewIHaveMgr(iHaveCh chan IHaveObj, entries *xn.IDMap,
	outMsgCh chan *UpaxClusterMsg, stopCh chan bool) (
	mgr *IHaveMgr, err error) {

	if iHaveCh == nil {
		err = NilIHaveChan
	} else if entries == nil {
		err = NilIDMap
	} else if outMsgCh == nil {
		err = NilOutMsgCh
	} else {
		mgr = &IHaveMgr{
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
func (mgr *IHaveMgr) Run() {
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
