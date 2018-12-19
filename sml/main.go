package sml // import "shanhu.io/sml/sml"

import (
	"shanhu.io/base/subcmd"
)

func cmd() *subcmd.List {
	c := subcmd.New()
	c.AddHost("sync", "synchronizes a profile", sync)
	c.AddHost("track", "tracks a repository in the profile", track)
	c.AddHost("untrack", "untracks a repository in the profile", untrack)
	c.SetDefaultServer("https://gopkgs.io")
	return c
}

// Main is the main entrance function for smlsync command.
func Main() {
	cmd().Main()
}
