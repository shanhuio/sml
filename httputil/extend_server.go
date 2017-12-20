package httputil

import (
	"net/url"
	"strings"
)

// ExtendServer extends the server host string with an https:// prefix if it is
// not a localhost, or an http:// prefix if it is localhost. It only extends
// the server when the server string can be extends to be a valid URL.
func ExtendServer(s string) string {
	if strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://") {
		return s
	}

	ret := "https://" + s
	u, err := url.Parse(ret)
	if err != nil {
		return s // not a valid URL
	}
	if u.Scheme != "https" {
		return s // not something that we think of
	}

	if u.Scheme == "https" {
		if u.Host == "localhost" || strings.HasPrefix(u.Host, "localhost:") {
			// testing on localhost, set the scheme to http
			u.Scheme = "http"
			return u.String()
		}
	}
	return ret
}
