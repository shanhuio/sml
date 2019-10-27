package core

// StateQuery is a state query.
type StateQuery struct {
	Tracking []string
}

// State is contains a state of a mapped super repo.
type State struct {
	ID      string
	Clock   uint64
	Commits map[string]string
	Sources map[string]string

	// Missing are repos requested but missing
	Missing []string `json:",omitempty"`
}

// NewState creates a new state with a particular ID.
func NewState(com *Commit) *State {
	return &State{
		ID:      com.ID,
		Clock:   com.Clock,
		Commits: make(map[string]string),
		Sources: make(map[string]string),
	}
}
