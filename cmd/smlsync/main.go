// Command smlsync synchronizes GOPATH repositories to forge.shanhu.io HEAD.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"shanhu.io/misc/jsonutil"
	"shanhu.io/sml/sync"
)

type runner struct {
	remote  *sync.RemoteState
	verbose bool
}

func knownHostsFile() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("get current user: %s", err)
	}
	return filepath.Join(u.HomeDir, ".shanhu/ssh_known_hosts"), nil
}

func (r *runner) run() error {
	state, err := r.remote.Fetch()
	if err != nil {
		return err
	}

	if r.verbose {
		jsonutil.Print(state)
	}

	knownHosts, err := knownHostsFile()
	if err != nil {
		return err
	}
	syncer := &sync.Syncer{
		KnownHostsFile: knownHosts,
	}
	return syncer.Sync(state, nil)
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

	remote := &sync.RemoteState{
		Server:  *server,
		Project: *proj,
	}
	if *mirror != "" {
		remote.Transform = sync.MirrorSourceTransform(*mirror)
	} else if *org != "" {
		remote.Transform = sync.BitbucketSourceTransform(*org)
	}

	r := &runner{
		remote:  remote,
		verbose: *verbose,
	}

	if err := r.run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
