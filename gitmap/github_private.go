package gitmap

import (
	"fmt"
	"strings"
)

// GitHubPrivate maps GitHub private repositories from https to git.
type GitHubPrivate struct {
	org string

	matchPrefix string
	addPrefix   string
}

// NewGitHubPrivate creates a new GitHub private mapper.
func NewGitHubPrivate(org string) *GitHubPrivate {
	return &GitHubPrivate{
		org: org,

		matchPrefix: fmt.Sprintf("https://github.com/%s/", org),
		addPrefix:   fmt.Sprintf("git@github.com:%s/", org),
	}
}

// Map maps a GitHub https source to git private source if matches the username
// or organization.
func (p *GitHubPrivate) Map(src string) string {
	if !strings.HasPrefix(src, p.matchPrefix) {
		return src
	}
	name := strings.TrimPrefix(src, p.matchPrefix)
	return p.addPrefix + name
}
