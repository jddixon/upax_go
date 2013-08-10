package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

func Usage() {
	fmt.Printf("Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Printf("where the options are:\n")
	flag.PrintDefaults()
}

const (
	DEFAULT_LFS = "/var/Upax"
)

var (
	// these need to be referenced as pointers
	justShow = flag.Bool("j", false, "display option settings and exit")
	lfs      = flag.String("lfs", DEFAULT_LFS, "path to work directory")
	testing  = flag.Bool("T", false, "test run")
	verbose  = flag.Int("v", 0, "how talkative to be")
)

func init() {
	fmt.Println("init() invocation") // DEBUG
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	// FIXUPS ///////////////////////////////////////////////////////
	if *testing {
		if *lfs == DEFAULT_LFS || lft == "" {
			*lfs = "./tmp"
		} else {
			*lfs = path.Join("tmp", *lfs)
		}
	}
	// SANITY CHECKS ////////////////////////////////////////////////

	// DISPLAY FLAGS ////////////////////////////////////////////////
	if *verbose > 0 || *justShow {
		fmt.Printf("justShow               = %v\n", *justShow)
		fmt.Printf("lfs                    = %s\n", *lfs)
		fmt.Printf("testing                = %v\n", *testing)
		fmt.Printf("verbose                = %d\n", *verbose)
	}
	if *justShow {
		return
	}
	// SET UP OPTIONS ///////////////////////////////////////////////

	// LOAD XLATTICE NODE CONFIGURATION /////////////////////////////

	// LOAD UPAX CLUSTER CONFIGURATION //////////////////////////////

	// LOAD FTLOG ///////////////////////////////////////////////////

	// SET UP COMMUNICATION WITH PEERS AND MIRRORS //////////////////

	// CATCH UP IF NECESSARY ////////////////////////////////////////

	// LISTEN AND SERVE /////////////////////////////////////////////

	// ORDERLY SHUTDOWN /////////////////////////////////////////////

}
