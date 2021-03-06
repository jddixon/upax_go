package upax_go

// upax_go/entry.go

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	xu "github.com/jddixon/xlUtil_go"
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
	key       []byte // SHA 1/2/3 content key
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
	} else if (len(key) != xu.SHA1_BIN_LEN && len(key) != xu.SHA3_BIN_LEN) ||
		(len(nodeID) != xu.SHA1_BIN_LEN && len(nodeID) != xu.SHA3_BIN_LEN) {
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
func (e *LogEntry) Timestamp() int64 {
	return e.timestamp
}

// Whether the key is an SHA1 key.
func (e *LogEntry) UsingSHA1() bool {
	return len(e.key) == xu.SHA1_BIN_LEN
}

// EQUAL ////////////////////////////////////////////////////////////

func (e *LogEntry) Equal(any interface{}) bool {
	if any == e {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *LogEntry:
		_ = v
	default:
		return false
	}
	other := any.(*LogEntry) // type assertion

	if e.timestamp != other.timestamp {
		return false
	} else if !bytes.Equal(e.key, other.key) {
		return false
	} else if !bytes.Equal(e.nodeID, other.nodeID) {
		return false
	} else if e.src != other.src {
		return false
	} else if e.path != other.path {
		return false
	}
	return true
}

// SERIALIZATION AND DESERIALIZATION ////////////////////////////////

func (e *LogEntry) String() string {
	return fmt.Sprintf("%d %s %s \"%s\" %s", e.timestamp,
		hex.EncodeToString(e.key), hex.EncodeToString(e.nodeID),
		e.src, e.path)
}

func ParseLogEntry(s string, whichSHA int) (e *LogEntry, err error) {
	var groups []string
	var t, key, nodeID, src, path string
	var tInt int
	var t64 int64
	var keyBuf, nodeIDBuf []byte

	s = strings.TrimSpace(s)
	switch whichSHA {
	case xu.USING_SHA1:
		groups = bodyLine1RE.FindStringSubmatch(s)
	case xu.USING_SHA2:
		groups = bodyLine2RE.FindStringSubmatch(s)
	case xu.USING_SHA3:
		groups = bodyLine3RE.FindStringSubmatch(s)
		// XXX DEFAULT = ERROR
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
