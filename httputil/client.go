package httputil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"log"
)

func makeURL(server, p string) (string, error) {
	u, err := url.Parse(server)
	if err != nil {
		return "", err
	}

	// append url
	u.Path = path.Join(u.Path, p)
	return u.String(), nil
}

// Client performs client that calls to a remote server with an optional token.
type Client struct {
	Server string
	Token  string // Optional token to be put in the Bearer HTTP header.
}

// NewClient creates a new client.
func NewClient(s string) *Client {
	return &Client{Server: s}
}

// NewTokenClient creates a new client with a Bearer token.
func NewTokenClient(s, tok string) *Client {
	return &Client{Server: s, Token: tok}
}

func (c *Client) addToken(req *http.Request) {
	if c.Token != "" {
		AddToken(req, c.Token)
	}
}

// Put puts a stream to a route on the server.
func (c *Client) Put(p string, r io.Reader) error {
	u, err := makeURL(c.Server, p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", u, r)
	if err != nil {
		return err
	}
	c.addToken(req)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return RespError(resp)
	}
	return nil
}

// Poke posts a signal to the given route on the server.
func (c *Client) Poke(p string) error {
	u, err := makeURL(c.Server, p)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", u, nil)
	if err != nil {
		return err
	}
	c.addToken(req)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return RespError(resp)
	}
	return nil
}

// Get gets a response from a route on the server.
func (c *Client) Get(p string) (*http.Response, error) {
	u, err := makeURL(c.Server, p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	c.addToken(req)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		return nil, RespError(resp)
	}

	return resp, nil
}

func (c *Client) jsonPost(p string, req interface{}) (*http.Response, error) {
	u, err := makeURL(c.Server, p)
	if err != nil {
		return nil, err
	}
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", u, bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	c.addToken(httpReq)

	httpClient := new(http.Client)
	return httpClient.Do(httpReq)
}

// JSONPost posts a JSON object as the request body and writes the body
// into the given writer.
func (c *Client) JSONPost(p string, req interface{}, w io.Writer) error {
	resp, err := c.jsonPost(p, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return RespError(resp)
	}

	if w == nil {
		return nil
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return nil
}

// JSONCall performs a call with the request as a marshalled JSON object,
// and the response unmarhsalled as a JSON object.
func (c *Client) JSONCall(p string, req, resp interface{}) error {
	httpResp, err := c.jsonPost(p, req)
	if err != nil {
		log.Println("json post", err)
		return err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return RespError(httpResp)
	}

	if resp == nil {
		return nil
	}
	dec := json.NewDecoder(httpResp.Body)
	return dec.Decode(resp)
}
