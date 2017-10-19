package core

import (
	"encoding/json"

	"shanhu.io/pisces/objects"
)

// Version is the version of smallrepo's super repo core.
const Version = "0.1"

// Core is the core of a build result.
type Core struct {
	Version   string
	GoVersion string
	Tracking  []string
	Commits   map[string]string
}

// Hash creates the hash of this core result object.
func (c *Core) Hash() string {
	bs, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return objects.Hash(bs)
}
