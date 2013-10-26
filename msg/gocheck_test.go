package msg

import (
	. "launchpad.net/gocheck"
	"testing"
)

// IF USING gocheck, need a file with next 3 lines in each package=directory.
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})
