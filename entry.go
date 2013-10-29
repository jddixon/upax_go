package upax_go

// upax_go/entry.go

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var (
	InvalidKeyOrNodeID = errors.New("invalid key or nodeID")
	NoDQuotesInSrc     = errors.New("src may not contain double quote")
	NotAValidLogEntry  = errors.New("not a well-formed log entry")
	NotAValidPath      = errors.New("not a well-formed path")
	NilKeyOrNodeID     = errors.New("nil key or nodeID")
)

type LogEntry struct {
	timestamp int64  // nanoseconds from the Epoch (1 Jan 1970 0:0:0)
	key       []byte // SHA 1 or 3 content key
	nodeID    []byte // of the committer
	src       string // free-form
	path      string // POSIX path, email address, and similar
}

func NewLogEntry(t int64, key []byte, nodeID []byte,
	src, path string) (e *LogEntry, err error) {

	if t == 0 {
		t = time.Now().UnixNano() // ns since epoch
	}
	if key == nil || nodeID == nil {
		err = NilKeyOrNodeID
	} else if (len(key) != 20 && len(key) != 32) ||
		(len(nodeID) != 20 && len(nodeID) != 32) {
		err = InvalidKeyOrNodeID
	}
	if err == nil {
		path = strings.TrimSpace(path)
		if !pathRE.MatchString(path) {
			err = NotAValidPath
		}
	}
	if err == nil {
		src = strings.TrimSpace(src)
		if strings.Contains(src, "\"") {
			err = NoDQuotesInSrc
		}
	}
	if err == nil {
		e = &LogEntry{t, key, nodeID, src, path}
	}
	return
}

// ATTRIBUTES ///////////////////////////////////////////////////////

// Return the key in hex.
func (e *LogEntry) Key() []byte {
	// XXX should be copy
	return e.key
}

// Return the nodeID in hex.
func (e *LogEntry) NodeID() []byte {
	// XXX should be copy
	return e.nodeID
}

// Return the path, which might be POSIX or an email address or ...
func (e *LogEntry) Path() string {
	return e.path
}

// Return the client-supplied source.
func (e *LogEntry) Src() string {
	return e.src
}

// Return the timestamp as nanoseconds from the epoch, 00:00 on 1 Jan 1970
func (e *LogEntry) TimeStamp() int64 {
	return e.timestamp
}

// Whether the key is an SHA1 key.
func (e *LogEntry) UsingSHA1() bool {
	return len(e.key) == 20
}

// SERIALIZATION AND DESERIALIZATION ////////////////////////////////

func (e *LogEntry) String() string {
	return fmt.Sprintf("%d %s %s \"%s\" %s", e.timestamp,
		hex.EncodeToString(e.key), hex.EncodeToString(e.nodeID),
		e.src, e.path)
}

func ParseLogEntry(s string, usingSHA1 bool) (e *LogEntry, err error) {
	var groups []string
	var t, key, nodeID, src, path string
	var tInt int
	var t64 int64
	var keyBuf, nodeIDBuf []byte

	s = strings.TrimSpace(s)
	if usingSHA1 {
		groups = bodyLine1RE.FindStringSubmatch(s)
	} else {
		groups = bodyLine3RE.FindStringSubmatch(s)
	}
	if groups == nil {
		err = NotAValidLogEntry
	} else {
		// all of these are strings
		t = groups[1]
		key = groups[2]
		nodeID = groups[3]
		src = groups[4]
		path = groups[5]

		tInt, err = strconv.Atoi(t)
		if err == nil {
			t64 = int64(tInt)
			keyBuf, err = hex.DecodeString(key)
		}
		if err == nil {
			nodeIDBuf, err = hex.DecodeString(nodeID)
		}
		if err == nil {
			e, err = NewLogEntry(t64, keyBuf, nodeIDBuf, src, path)
		}
	}
	return
}
