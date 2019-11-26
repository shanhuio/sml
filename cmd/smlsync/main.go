// Command smlsync synchronizes GOPATH repositories to forge.shanhu.io HEAD.
package main

import (
	"flag"
	"fmt"
	"os"

	"shanhu.io/aries/creds"
	"shanhu.io/misc/jsonutil"
	"shanhu.io/sml/core"
	"shanhu.io/sml/gitmap"
	"shanhu.io/sml/sync"
)

func run(server, proj, org string, verbose bool) error {
	c, err := creds.Dial(server)
	if err != nil {
		return err
	}

	state := new(core.State)
	if err := c.JSONCall("/api/sync/proj", proj, state); err != nil {
		return err
	}

	gitmap.MapCoreState(state, gitmap.NewBitbucketPrivate(org))

	if verbose {
		jsonutil.Print(state)
	}

	return sync.Sync(nil, state, nil)
}

func main() {
	server := flag.String(
		"server", "https://forge.shanhu.io", "Server address.",
	)
	org := flag.String("org", "shanhuio", "Default private org on bitbucket.")
	proj := flag.String("proj", "h8liu", "Project to sync to.")
	verbose := flag.Bool("v", false, "print the state")
	flag.Parse()

	if err := run(*server, *proj, *org, *verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
