package msfrpc

import "fmt"

type sessionListReq struct {
	_msgpack struct{} `msgpack:",asArray"`
	Method   string
	Token    string
}

type SessionListResp struct {
	ID          int    `msgpack:",omitempty"`
	Type        string `msgpack:"type"`
	TunnelLocal string `msgpack:"tunnel_local"`
	TunnelPeer  string `msgpack:"tunnel_peer"`
	ViaExploit  string `msgpack:"via_exploit"`
	ViaPayload  string `msgpack:"via_payload"`
	Description string `msgpack:"desc"`
	Info        string `msgpack:"info"`
	Workspace   string `msgpack:"workspace"`
	SessionHost string `msgpack:"session_host"`
	SessionPort int    `msgpack:"session_port"`
	Username    string `msgpack:"username"`
	UUID        string `msgpack:"uuid"`
	ExploitUUID string `msgpack:"exploit_uuid"`
}

func (c Client) SessionList() (map[int]SessionListResp, error) {
	req := sessionListReq{
		Method: "session.list",
		Token:  c.token,
	}

	resp := make(map[int]SessionListResp)
	err := c.Send(req, &resp)
	if err != nil {
		return nil, fmt.Errorf("error getting sessionlist: %v", err)
	}

	// flatten map
	for id, session := range resp {
		session.ID = id
		resp[id] = session
	}
	return resp, nil
}
