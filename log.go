package upax_go

// upax_go/log.go

import ()

type Log struct {
	entries    *LogEntry
	usingSHA1  bool
	timestamp  uint64 // of creation of the current log
	master     []byte // creator of this log
	prevHash   []byte // hash of previous log
	prevMaster []byte // nodeID, owner of previous log
}
