// Package gitmap provides mapping functions to convert https git sources to
// SSH git sources.
package gitmap

import (
	"fmt"
	"strings"
)

// BitbucketPrivate maps bitbucket private repositories from https to git.
type BitbucketPrivate struct {
	org string

	matchPrefix string
	addPrefix   string
}

// NewBitbucketPrivate creates a new Bitbucket private mapper.
func NewBitbucketPrivate(org string) *BitbucketPrivate {
	return &BitbucketPrivate{
		org: org,

		matchPrefix: fmt.Sprintf("https://bitbucket.org/%s/", org),
		addPrefix:   fmt.Sprintf("git@bitbucket.org:%s/", org),
	}
}

// Map maps a bitbucket https source to git private source if it matches
// organization.
func (p *BitbucketPrivate) Map(src string) string {
	if !strings.HasPrefix(src, p.matchPrefix) {
		return src
	}
	name := strings.TrimPrefix(src, p.matchPrefix)
	return p.addPrefix + name
}
