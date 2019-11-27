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

type runner struct {
	server   string
	proj     string
	stateMap func(s *core.State)
	verbose  bool
}

func (r *runner) run() error {
	c, err := creds.Dial(r.server)
	if err != nil {
		return err
	}

	state := new(core.State)
	if err := c.JSONCall("/api/sync/proj", r.proj, state); err != nil {
		return err
	}

	if r.stateMap != nil {
		r.stateMap(state)
	}

	if r.verbose {
		jsonutil.Print(state)
	}

	return sync.Sync(nil, state, nil)
}

func main() {
	server := flag.String(
		"server", "https://forge.shanhu.io", "Server address.",
	)
	org := flag.String("org", "shanhuio", "Default private org on bitbucket.")
	mirror := flag.String("mirror", "", "Sync from this mirror machine.")
	proj := flag.String("proj", "h8liu", "Project to sync to.")
	verbose := flag.Bool("v", false, "print the state")
	flag.Parse()

	r := &runner{
		server:  *server,
		proj:    *proj,
		verbose: *verbose,
	}

	if *mirror != "" {
		r.stateMap = func(s *core.State) {
			srcs := make(map[string]string)
			for repo := range s.Sources {
				srcs[repo] = *mirror + "/" + repo
			}
			s.Sources = srcs
		}
	} else if *org != "" {
		r.stateMap = func(s *core.State) {
			gitmap.MapCoreState(s, gitmap.NewBitbucketPrivate(*org))
		}
	}

	if err := r.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
