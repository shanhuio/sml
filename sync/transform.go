package sync

import (
	"shanhu.io/sml/core"
	"shanhu.io/sml/gitmap"
)

// MirrorSourceTransform maps all sources to a mirror prefix.  The git source
// is formatted as (mirror + "/" + <repo name>).
func MirrorSourceTransform(mirror string) func(s *core.State) {
	return func(s *core.State) {
		srcs := make(map[string]string)
		for repo := range s.Sources {
			srcs[repo] = mirror + "/" + repo
		}
		s.Sources = srcs
	}
}

// PrivateSourceTransform maps repo soruces in a Bitbucket and GitHub org from
// https to git+ssh format.
func PrivateSourceTransform(org string) func(s *core.State) {
	return func(s *core.State) {
		gitmap.MapCoreState(s, gitmap.NewBitbucketPrivate(org))
		gitmap.MapCoreState(s, gitmap.NewGitHubPrivate(org))
	}
}
