package gitmap

import (
	"shanhu.io/sml/core"
)

// Mapper maps a git source.
type Mapper interface {
	Map(src string) string
}

// MapCoreState maps all the git source in state.
func MapCoreState(state *core.State, m Mapper) {
	srcs := make(map[string]string)
	for repo, src := range state.Sources {
		srcs[repo] = m.Map(src)
	}
	state.Sources = srcs
}
