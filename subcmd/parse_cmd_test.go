package subcmd

import (
	"testing"
)

func TestParseCmd(t *testing.T) {
	for _, test := range []struct {
		in, name, host string
	}{
		{"", "", ""},
		{"hello", "hello", ""},
		{"hello@something", "hello", "something"},
		{"@something", "", "something"},
		{"a@b@c", "a", "b@c"},
		{"hello@", "hello", ""},
	} {
		name, host := parseCmd(test.in)
		if name != test.name || host != test.host {
			t.Errorf(
				"parseCmd(%q), want (%q, %q), got (%q, %q)",
				test.in, test.name, test.host,
				name, host,
			)
		}
	}
}
