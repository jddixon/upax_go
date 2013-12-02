package ${pkgName}

// ${pkgName}/${filePrefix}ihave_mgr.go

import (
	xn "github.com/jddixon/xlattice_go/node"
)

// The ${TypePrefix}IHaveMgr receives IHaveObjs on its input channel.  For each
// ID, it checks to see whether the ID is present in its entries.
// If the ID is not present, it creates a Get message, which it passes
// on.
//
type ${TypePrefix}IHaveMgr struct {
	iHaveCh  chan IHaveObj
	entries  *xn.IDMap // a convenience
	outMsgCh chan *Upax${TypePrefix}Msg
	stopCh   chan bool
}

func New${TypePrefix}IHaveMgr(iHaveCh chan IHaveObj, entries *xn.IDMap,
	outMsgCh chan *Upax${TypePrefix}Msg, stopCh chan bool) (
	mgr *${TypePrefix}IHaveMgr, err error) {

	if iHaveCh == nil {
		err = NilIHaveChan
	} else if entries == nil {
		err = NilIDMap
	} else if outMsgCh == nil {
		err = NilOutMsgCh
	} else {
		mgr = &${TypePrefix}IHaveMgr{
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
func (mgr *${TypePrefix}IHaveMgr) Run() {
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
					op := Upax${TypePrefix}Msg_Get
					msgOut := &Upax${TypePrefix}Msg{
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
