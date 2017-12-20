package httputil

import (
	"net/http"
)

// AddToken adds the authorization header into the request.
func AddToken(req *http.Request, tok string) {
	req.Header.Set("Authorization", "Bearer "+tok)
}
