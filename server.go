package upax_go

// upax_go/server.go

import (
	"crypto/rsa"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
)

var _ = fmt.Print

type UpaxServer struct {
	ClusterName, ServerName string
	ClusterID               *xi.NodeID
	ckPriv, skPriv          *rsa.PrivateKey
	xn.Node
}

func NewUpaxServer(clusterName, serverName string, clusterID *xi.NodeID,
	ckPriv, skPriv *rsa.PrivateKey, node *xn.Node) (us *UpaxServer, err error) {

	if clusterName == "" {
		clusterName = "upax"
	}
	if serverName == "" {
		err = EmptyName
	} else if ckPriv == nil || ckPriv == nil {
		err = NilRSAKey
	} else if node == nil {
		err = NilNode
	}
	if err == nil {
		us = &UpaxServer{
			ClusterName: clusterName,
			ServerName:  serverName,
			ckPriv:      ckPriv,
			skPriv:      skPriv,
			Node:        *node,
		}
	}
	return
}
