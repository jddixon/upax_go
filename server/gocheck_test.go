package server

import (
	. "launchpad.net/gocheck"
	"testing"
)

// IF USING gocheck, need a file with the next 3 lines in each directory.
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})
