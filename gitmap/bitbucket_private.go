// Package gitmap provides mapping functions to convert https git sources to
// SSH git sources.
package gitmap

import (
	"fmt"
	"strings"
)

// BitbucketPrivate maps bitbucket private buckets from https to git.
type BitbucketPrivate struct {
	org    string
	prefix string
}

// NewBitbucketPrivate creates a new Bitbucket private mapper.
func NewBitbucketPrivate(org string) *BitbucketPrivate {
	return &BitbucketPrivate{
		org:    org,
		prefix: fmt.Sprintf("https://bitbucket.org/%s/", org),
	}
}

// Map maps a bitbucket https source to git private source.
func (p *BitbucketPrivate) Map(src string) string {
	if !strings.HasPrefix(src, p.prefix) {
		return src
	}
	name := strings.TrimPrefix(src, p.prefix)
	return fmt.Sprintf("git@bitbucket.org:%s/%s", p.org, name)
}
