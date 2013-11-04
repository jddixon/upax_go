package upax_go

// upax_go/server.go

import (
	"crypto/rsa"
	"fmt"
	"github.com/jddixon/xlattice_go/reg"
	"github.com/jddixon/xlattice_go/u"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	"path/filepath"
)

var _ = fmt.Print

type UpaxServer struct {
	uDir           u.UI
	ckPriv, skPriv *rsa.PrivateKey
	reg.ClusterMember
}

func NewUpaxServer(ckPriv, skPriv *rsa.PrivateKey, cm *reg.ClusterMember) (
	us *UpaxServer, err error) {

	var (
		lfs  string
		uDir u.UI
	)
	if ckPriv == nil || ckPriv == nil {
		err = NilRSAKey
	} else if cm == nil {
		err = NilClusterMember
	}
	if err == nil {
		// whatever created cm should have created the local file system
		// and written the node configuration to
		// LFS/.xlattice/cluster.member.config.  Let's make sure that
		// that exists before proceeding.

		lfs = cm.GetLFS()
		pathToCfg := filepath.Join(
			filepath.Join(lfs, ".xlattice"), "cluster.member.config")
		var found bool
		found, err = xf.PathExists(pathToCfg)
		if err == nil && found == false {
			err = ClusterConfigNotFound
		}
	}
	if err == nil {
		// DEBUG
		fmt.Printf("creating directory tree in %s\n", lfs)
		// END

		pathToU := filepath.Join(lfs, "U")
		uDir, err = u.New(pathToU, u.DIR16x16, 0)
	}
	if err == nil {
		us = &UpaxServer{
			uDir:          uDir,
			ckPriv:        ckPriv,
			skPriv:        skPriv,
			ClusterMember: *cm,
		}
	}
	return
}

func (us *UpaxServer) Run() {

	// XXX STUB XXX

}
