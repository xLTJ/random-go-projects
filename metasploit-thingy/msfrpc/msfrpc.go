package msfrpc

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"github.com/vmihailenco/msgpack/v5"
	"net/http"
)

// Docs: https://usermanual.wiki/Pdf/metasploitrpcguide.221304445.pdf

type Client struct {
	host     string
	port     string
	username string
	password string
	token    string
}

type errorResp struct {
	Error        bool   `msgpack:"error"`
	ErrorClass   string `msgpack:"error_class"`
	ErrorMessage string `msgpack:"error_message"`
}

func NewClient() (Client, error) {
	host := viper.GetString("msgrpc.host")
	port := viper.GetString("msgrpc.port")
	username := viper.GetString("msgrpc.username")
	password := viper.GetString("msgrpc.password")

	if len(host) == 0 || len(port) == 0 || len(username) == 0 || len(password) == 0 {
		return Client{}, fmt.Errorf("missing one or more values. make sure the following are set:\n    host, port, username, password\n")
	}

	client := Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}

	token, err := client.Login()
	if err != nil {
		return Client{}, err
	}

	client.token = token
	return client, nil
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

	dest := fmt.Sprintf("http://%s:%s/api", c.host, c.port)
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
