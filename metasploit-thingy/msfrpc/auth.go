package msfrpc

import "fmt"

type loginReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Username string
	Password string
}

type loginResp struct {
	Result     string `msgpack:"result"`
	Token      string `msgpack:"token"`
	Error      bool   `msgpack:"error"`
	ErrorClass string `msgpack:"error_class"`
}

type logoutReq struct {
	_msgpack    struct{} `msgpack:",asArray"`
	Method      string
	Token       string
	LogoutToken string
}

type logoutResp struct {
	Result string `msgpack:"result"`
}

func (c Client) Login() (string, error) {
	req := loginReq{
		Method:   "auth.login",
		Username: c.username,
		Password: c.password,
	}

	var resp loginResp
	err := c.Send(req, &resp)
	if err != nil {
		return "", fmt.Errorf("Error sending login request: %v\n", err)
	}

	if resp.Error {
		return "", fmt.Errorf("Error logging in (probably invalid credentials): %v\n", resp.ErrorClass)
	}

	return resp.Token, nil
}

func (c Client) Logout() error {
	req := logoutReq{
		Method:      "auth.logout",
		Token:       c.token,
		LogoutToken: c.token,
	}

	var resp logoutResp
	err := c.Send(req, &resp)
	if err != nil {
		return err
	}

	c.token = ""
	return nil
}
