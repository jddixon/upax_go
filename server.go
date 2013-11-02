package upax_go

// upax_go/server.go

import (
	"crypto/rsa"
	"fmt"
	"github.com/jddixon/xlattice_go/reg"
)

var _ = fmt.Print

type UpaxServer struct {
	ckPriv, skPriv *rsa.PrivateKey
	reg.ClusterMember
}

func NewUpaxServer(ckPriv, skPriv *rsa.PrivateKey, cm *reg.ClusterMember) (
	us *UpaxServer, err error) {

	if ckPriv == nil || ckPriv == nil {
		err = NilRSAKey
	} else if cm == nil {
		err = NilClusterMember
	}
	if err == nil {
		us = &UpaxServer{
			ckPriv:        ckPriv,
			skPriv:        skPriv,
			ClusterMember: *cm,
		}
	}
	return
}
