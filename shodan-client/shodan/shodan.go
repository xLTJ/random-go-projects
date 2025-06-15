package shodan

const ApiBaseUrl = "https://api.shodan.io"

type Client struct {
	apiKey string
}

func (c Client) GetApiKey() string {
	return c.apiKey
}

func NewClient(apiKey string) Client {
	return Client{apiKey: apiKey}
}
