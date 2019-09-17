package core

import (
	"shanhu.io/misc/strutil"
)

// Change runes in DiffLine
const (
	DiffAdd     = '+'
	DiffRemove  = '-'
	DiffReplace = ' '
)

// DiffLine is one single difference between two commit cores.
type DiffLine struct {
	Change rune
	From   string
	To     string
}

func makeDiffline(from, to string) *DiffLine {
	if from == to {
		return nil
	}
	if from == "" {
		return &DiffLine{
			Change: DiffAdd,
			To:     to,
		}
	}
	if to == "" {
		return &DiffLine{
			Change: DiffRemove,
			From:   from,
		}
	}
	return &DiffLine{
		Change: DiffReplace,
		From:   from,
		To:     to,
	}
}

// Diff contains a diff between two commit cores.
type Diff struct {
	Version   *DiffLine
	GoVersion *DiffLine
	Tracking  []*DiffLine
	Commits   map[string]*DiffLine

	MultiParent bool
}

// NewDiff creates a new diff for two cores.
func NewDiff(from, to *Core) *Diff {
	if from == nil {
		ret := &Diff{
			Version:   makeDiffline("", to.Version),
			GoVersion: makeDiffline("", to.GoVersion),
			Commits:   make(map[string]*DiffLine),
		}

		for _, repo := range to.Tracking {
			d := makeDiffline("", repo)
			ret.Tracking = append(ret.Tracking, d)
		}

		for repo, commit := range to.Commits {
			ret.Commits[repo] = makeDiffline("", commit)
		}
		return ret
	}

	ret := &Diff{
		Version:   makeDiffline(from.Version, to.Version),
		GoVersion: makeDiffline(from.GoVersion, to.GoVersion),
		Commits:   make(map[string]*DiffLine),
	}

	setTo := strutil.MakeSet(to.Tracking)
	for _, repo := range from.Tracking {
		if !setTo[repo] {
			d := makeDiffline(repo, "")
			ret.Tracking = append(ret.Tracking, d)
		}
	}

	setFrom := strutil.MakeSet(from.Tracking)
	for _, repo := range to.Tracking {
		if !setFrom[repo] {
			d := makeDiffline("", repo)
			ret.Tracking = append(ret.Tracking, d)
		}
	}

	for repo, commit := range from.Commits {
		d := makeDiffline(commit, to.Commits[repo])
		if d != nil {
			ret.Commits[repo] = d
		}
	}
	for repo, commit := range to.Commits {
		_, found := from.Commits[repo]
		if !found {
			ret.Commits[repo] = makeDiffline("", commit)
		}
	}

	return ret
}
