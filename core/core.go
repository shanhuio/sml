package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

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

	ret := sha256.Sum256(bs)
	return hex.EncodeToString(ret[:])
}
