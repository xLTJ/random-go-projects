package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiInfo struct {
	ScanCredits  int    `json:"scan_credits"`
	Plan         string `json:"plan"`
	Http         bool   `json:"http"`
	Unlocked     bool   `json:"unlocked"`
	QueryCredits int    `json:"query_credits"`
	MonitoredIps int    `json:"monitored_ips"`
	UnlockedLeft int    `json:"unlocked_left"`
	Telnet       bool   `json:"telnet"`
}

// ApiInfo retrieves information about the API key's usage limits, plan details, and available credits.
func (c Client) ApiInfo() (ApiInfo, error) {
	res, err := http.Get(fmt.Sprintf("%s/api-info?key=%s", ApiBaseUrl, c.apiKey))
	if err != nil {
		return ApiInfo{}, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		return ApiInfo{}, fmt.Errorf("unable to call api: %s", res.Status)
	}

	var apiInfo ApiInfo
	err = json.NewDecoder(res.Body).Decode(&apiInfo)
	if err != nil {
		return ApiInfo{}, err
	}
	return apiInfo, nil
}
