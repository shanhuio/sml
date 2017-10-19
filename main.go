package sml

import (
	"shanhu.io/misc/subcmd"
)

func cmd() *subcmd.List {
	c := subcmd.New()
	c.AddHost("sync", "synchronizes a profile", sync)
	c.AddHost("track", "tracks a repository in the profile", track)
	c.AddHost("untrack", "untracks a repository in the profile", untrack)
	c.SetDefaultServer("https://smallrepo.com")
	return c
}

// Main is the main entrance function for smlsync command.
func Main() { cmd().Main() }
