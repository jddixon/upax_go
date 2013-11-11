package upax_go

import (
	e "errors"
)

var (
	BadDigSig             = e.New("bad digital signature")
	ClusterConfigNotFound = e.New("cluster config not found")
	EmptyName             = e.New("empty name parameter")
	NilClusterMember      = e.New("nil cluster member parameter")
	NilID                 = e.New("nil ID parameter")
	NilNode               = e.New("nil node parameter")
	NilRSAKey             = e.New("nil RSA private key parameter")
	NilServer             = e.New("nil UpaxServer parameter")
	NilUDir               = e.New("nil uDir parameter")
	NoDigSig              = e.New("no digital signature found")
	NoMembers             = e.New("nil or empty members parameter")
	NoSigFields           = e.New("no dig sig data fields found")
	NotClusterMember      = e.New("not a cluster member")
	UnknownClient         = e.New("unknown client")
)
