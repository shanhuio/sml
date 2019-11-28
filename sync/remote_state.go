package sync

import (
	"shanhu.io/aries/creds"
	"shanhu.io/sml/core"
)

func fetchCoreState(server, proj string) (*core.State, error) {
	c, err := creds.Dial(server)
	if err != nil {
		return nil, err
	}

	state := new(core.State)
	if err := c.JSONCall("/api/sync/proj", proj, state); err != nil {
		return nil, err
	}

	return state, nil
}

// RemoteState fetches a state from a remote /api/sync/proj URL.
type RemoteState struct {
	Server    string
	Project   string
	Transform func(s *core.State)
}

// Fetch fetches the state.
func (s *RemoteState) Fetch() (*core.State, error) {
	state, err := fetchCoreState(s.Server, s.Project)
	if err != nil {
		return nil, err
	}

	if s.Transform != nil {
		s.Transform(state)
	}
	return state, nil
}
