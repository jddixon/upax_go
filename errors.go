package upax_go

import (
	e "errors"
)

var (
	EmptyName = e.New("empty name parameter")
	NilID     = e.New("nil ID parameter")
	NilNode   = e.New("nil node parameter")
	NilRSAKey = e.New("nil RSA private key parameter")
)
