package core

import (
	"time"
)

// Commit is a commit of the super repo.
type Commit struct {
	ID      string    // a random id
	Time    time.Time // commit time
	Clock   uint64    // logical time
	Parents []string  // parent ids
	Super   string    // id to the super commit
	Message []byte    // a brief message
	Data    string    // hash to core data
}

// CommitCore is a commit of the super repo and its payload -- the core.
type CommitCore struct {
	Commit *Commit
	Core   *Core
}
