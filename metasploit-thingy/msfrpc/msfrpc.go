package msfrpc

import (
	"bytes"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"net/http"
)

// Docs: https://usermanual.wiki/Pdf/metasploitrpcguide.221304445.pdf

type Client struct {
	host     string
	username string
	password string
	token    string
}

func NewClient(host, username, password, token string) Client {
	return Client{
		host:     host,
		username: username,
		password: password,
	}
}

// Send sends a request in messagepack form, and returns the response
// using pointer based response return stuff because we need decode to know the form of the response
// also cus u cant use generics with methods, ughhh
func (c Client) Send(req, resp interface{}) error {
	buf := new(bytes.Buffer)
	err := msgpack.NewEncoder(buf).Encode(req)
	if err != nil {
		return err
	}

	dest := fmt.Sprintf("http://%s/api", c.host)
	httpResp, err := http.Post(dest, "binary/message-pack", buf)
	if err != nil {
		return err
	}
	defer func() { _ = httpResp.Body.Close() }()

	if err = msgpack.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return err
	}
	return nil
}
