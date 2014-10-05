package upax_go

// upax_go/memory.go

import (
	"fmt"
	"os"
)
var _ = fmt.Print

/**
 * What arrives on the wire.
 */
type Datum struct {
	SeqNbr		uint64
	Value		[]byte		// length of SHA1 hash
	RunningXOR	[]byte		// length of SHA1 hash
}

/**
 * This very simple state machine stores the value of each content
 * key in the order received.  
 */
type Memory struct {
	Values		[][]byte
	RunningXOR	[]byte		// of Memory.Values
	Pending		[]*Datum
	DataFile	*os.File	// open in append mode
	UsingSHA1	bool
}

func NewMemory(pathToDataFile string, usingSHA1 bool) (m *Memory, err error) {

	// XXX open the data file in binary append mode; we will want this to
	// be flushed on every write with a call to f.Sync() and eventually 
	// closed with f.Close().
	f, err := os.OpenFile(pathToDataFile, os.O_CREATE | os.O_APPEND, 0666)
	if err == nil {
		m = &Memory{
			DataFile:	f,
			UsingSHA1:	usingSHA1,
		}
	}
	return
}

func (m *Memory) Add(d *Datum) (err error) {

	if d.SeqNbr == m.NextSeqNbr() {
		m.Values = append(m.Values, d.Value)
	} else {
		m.Pending = append(m.Pending, d)
	}
	return
}
func (m *Memory) NextSeqNbr() uint64 {
	return uint64(len(m.Values))
}


