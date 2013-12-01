package upax_go

import (
	e "errors"
)

var (
	BadDigSig              = e.New("bad digital signature")
	BadHashLength          = e.New("bad SHA hash length")
	BadHash                = e.New("content hash doesn't match key")
	ClusterConfigNotFound  = e.New("cluster config not found")
	EmptyName              = e.New("empty name parameter")
	IntervalMustBePositive = e.New("interval must be positive")
	MissingTokenField      = e.New("missing token field")
	NilClusterMember       = e.New("nil cluster member parameter")
	NilHash                = e.New("nil hash field in message")
	NilID                  = e.New("nil ID parameter")
	NilIDMap               = e.New("nil IDMap parameter")
	NilIHaveChan           = e.New("nil IHaveCh parameter")
	NilLogEntry            = e.New("nil log entry field in message")
	NilMsgCh               = e.New("nil msgCh parameter")
	NilNode                = e.New("nil node parameter")
	NilOutMsgCh            = e.New("nil OutMsgCh parameter")
	NilPayload             = e.New("nil payload in message")
	NilRSAKey              = e.New("nil RSA private key parameter")
	NilServer              = e.New("nil UpaxServer parameter")
	NilToken               = e.New("nil token in message")
	NilUDir                = e.New("nil uDir parameter")
	NoDigSig               = e.New("no digital signature found")
	NoMembers              = e.New("nil or empty members parameter")
	NoSigFields            = e.New("no dig sig data fields found")
	UnknownPeer            = e.New("unknown peer")
)
