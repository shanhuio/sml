package httputil

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

// RespError returns the error from an HTTP response.
func RespError(resp *http.Response) error {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return errors.New(resp.Status + " - " + strings.TrimSpace(string(bs)))
}
