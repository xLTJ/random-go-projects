package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type IPHost struct {
	Tags        []string     `json:"tags"`
	IP          int          `json:"ip"`
	IPString    string       `json:"ip_str"`
	Domains     []string     `json:"domains"`
	CountryName string       `json:"country_name"`
	City        string       `json:"city"`
	OS          string       `json:"os"`
	Ports       []int        `json:"ports"`
	ISP         string       `json:"isp"`
	Org         string       `json:"org"`
	Data        []IPHostData `json:"data"`
}

type IPHostData struct {
	Port      int                  `json:"port"`
	Transport string               `json:"transport"`
	Version   string               `json:"version"`
	Product   string               `json:"product"`
	Vulns     map[string]HostVulns `json:"vulns"`
}

type HostVulns struct {
	Verified bool    `json:"verified"`
	CVSS     float64 `json:"cvss"`
}

func (c Client) IPHost(IPAddress string) (IPHost, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/%s?key=%s", ApiBaseUrl, IPAddress, c.apiKey))
	if err != nil {
		return IPHost{}, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		return IPHost{}, fmt.Errorf("unable to call api: %s", res.Status)
	}

	var ipHost IPHost
	err = json.NewDecoder(res.Body).Decode(&ipHost)
	if err != nil {
		return IPHost{}, err
	}
	return ipHost, nil
}
