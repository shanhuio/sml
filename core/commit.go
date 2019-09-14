package core

import (
	"time"
)

// Commit is a commit of the super repo.
type Commit struct {
	ID      string    // a random ID
	Project string    `json:",omitempty"`
	Time    time.Time // commit time
	Clock   uint64    // logical time
	Parents []string  // parent ids

	// Super (deprecated) is the ID to the super commit, the commit where this
	// commit is based on.
	//
	// Super string

	Message []byte // a brief message
	Data    string // hash to core data
}

// CommitCore is a commit of the super repo and its payload -- the core.
type CommitCore struct {
	Commit *Commit
	Core   *Core
}
