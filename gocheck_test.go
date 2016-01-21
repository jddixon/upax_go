package upax_go

import (
	. "gopkg.in/check.v1"
	"testing"
)

// IF USING gocheck, need a file with the next 3 lines in each directory.
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

const (
	VERBOSITY = 1
)
