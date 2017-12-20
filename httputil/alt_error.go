package httputil

import (
	"log"
	"net/http"
)

// AltError logs the actual error and replies with a message and code.
func AltError(w http.ResponseWriter, err error, msg string, code int) {
	log.Println(err)
	http.Error(w, msg, code)
}
