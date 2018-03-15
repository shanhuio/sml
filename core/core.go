package core

import (
	"encoding/json"

	"smallrepo.com/base/hashutil"
)

// Core is the core of a build result.
type Core struct {
	Version     string
	GoVersion   string
	Tracking    []string
	Commits     map[string]string
	BuildConfig *BuildConfig
}

// Hash creates the hash of this core result object.
func (c *Core) Hash() string {
	bs, err := json.Marshal(c)
	if err != nil {
		panic(err)
	}

	return hashutil.Hash(bs)
}
