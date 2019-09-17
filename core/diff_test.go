package core

import (
	"testing"

	"reflect"
)

func cleanDiff(d *Diff) {
	if len(d.Tracking) == 0 {
		d.Tracking = nil
	}
	if len(d.Commits) == 0 {
		d.Commits = nil
	}
}

func TestDiff(t *testing.T) {
	add := func(to string) *DiffLine {
		return &DiffLine{Change: DiffAdd, To: to}
	}

	sub := func(from string) *DiffLine {
		return &DiffLine{Change: DiffRemove, From: from}
	}

	change := func(from, to string) *DiffLine {
		return &DiffLine{
			Change: DiffReplace,
			From:   from,
			To:     to,
		}
	}

	for _, test := range []struct {
		from, to *Core
		diff     *Diff
	}{
		{
			from: &Core{},
			to:   &Core{},
			diff: &Diff{},
		}, {
			from: &Core{},
			to:   &Core{Version: "new"},
			diff: &Diff{
				Version: add("new"),
			},
		}, {
			from: &Core{},
			to:   &Core{GoVersion: "new go"},
			diff: &Diff{
				GoVersion: add("new go"),
			},
		}, {
			from: &Core{},
			to:   &Core{Tracking: []string{"shanhu.io/smlvm"}},
			diff: &Diff{
				Tracking: []*DiffLine{{
					Change: '+',
					To:     "shanhu.io/smlvm",
				}},
			},
		}, {
			from: &Core{Tracking: []string{"shanhu.io/smlvm"}},
			to:   &Core{},
			diff: &Diff{
				Tracking: []*DiffLine{sub("shanhu.io/smlvm")},
			},
		}, {
			from: &Core{Tracking: []string{"shanhu.io/smlvm"}},
			to:   &Core{Tracking: []string{"shanhu.io/smlvm"}},
			diff: &Diff{},
		}, {
			from: &Core{Tracking: []string{"shanhu.io/smlvm"}},
			to:   &Core{Tracking: []string{"shanhu.io/tools"}},
			diff: &Diff{
				Tracking: []*DiffLine{
					sub("shanhu.io/smlvm"),
					add("shanhu.io/tools"),
				},
			},
		}, {
			from: &Core{},
			to: &Core{Commits: map[string]string{
				"shanhu.io/smlvm": "abcdefg",
			}},
			diff: &Diff{
				Commits: map[string]*DiffLine{
					"shanhu.io/smlvm": add("abcdefg"),
				},
			},
		}, {
			from: &Core{Commits: map[string]string{
				"shanhu.io/smlvm": "abcdefg",
			}},
			to: &Core{},
			diff: &Diff{
				Commits: map[string]*DiffLine{
					"shanhu.io/smlvm": sub("abcdefg"),
				},
			},
		}, {
			from: &Core{Commits: map[string]string{
				"shanhu.io/smlvm": "abcdefg",
			}},
			to: &Core{Commits: map[string]string{
				"shanhu.io/smlvm": "cdefgab",
			}},
			diff: &Diff{
				Commits: map[string]*DiffLine{
					"shanhu.io/smlvm": change("abcdefg", "cdefgab"),
				},
			},
		},
	} {
		got := NewDiff(test.from, test.to)
		cleanDiff(test.diff)
		cleanDiff(got)
		if !reflect.DeepEqual(got, test.diff) {
			t.Errorf(
				"diff %v -> %v: got %v, want %v",
				test.from, test.to,
				got, test.diff,
			)
		}
	}
}
