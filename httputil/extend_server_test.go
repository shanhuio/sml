package httputil

import (
	"testing"
)

func TestExtendServer(t *testing.T) {
	o := func(s, want string) {
		got := ExtendServer(s)
		if got != want {
			t.Errorf("extendServer(%q), got %q want %q", s, got, want)
		}
	}

	o("http://localhost", "http://localhost")
	o("https://shanhu.io", "https://shanhu.io")
	o("localhost", "http://localhost")
	o("localhost:3356", "http://localhost:3356")
	o("shanhu.io", "https://shanhu.io")
}
