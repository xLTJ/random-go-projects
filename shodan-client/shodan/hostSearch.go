package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HostLocation struct {
	City         string  `json:"city"`
	RegionCode   string  `json:"region_code"`
	AreaCode     int     `json:"area_code"`
	Longitude    float64 `json:"longitude"`
	CountryCode3 string  `json:"country_code3"`
	Latitude     float64 `json:"latitude"`
	PostalCode   string  `json:"postal_code"`
	DmaCode      int     `json:"dma_code"`
	CountryCode  string  `json:"country_code"`
	CountryName  string  `json:"country_name"`
}

type Host struct {
	Product   string       `json:"product"`
	IP        uint32       `json:"ip"`
	Org       string       `json:"org"`
	ISP       string       `json:"isp"`
	Data      string       `json:"data"`
	ASN       string       `json:"asn"`
	Port      int          `json:"port"`
	Hostnames []string     `json:"hostnames"`
	Location  HostLocation `json:"location"`
	Timestamp string       `json:"timestamp"`
	Domains   []string     `json:"domains"`
	OS        string       `json:"os"`
	IPString  string       `json:"ip_str"`
}

type HostSearch struct {
	Matches []Host `json:"matches"`
}

// HostSearch sends a query to the Shodan API to search for hosts and returns the search results or an error.
func (c Client) HostSearch(query string) (HostSearch, error) {
	res, err := http.Get(fmt.Sprintf("%s/shodan/host/search?key=%s&query=%s", ApiBaseUrl, c.apiKey, query))
	if err != nil {
		return HostSearch{}, err
	}
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != 200 {
		return HostSearch{}, fmt.Errorf("unable to call api: %s", res.Status)
	}

	var hostSearch HostSearch
	err = json.NewDecoder(res.Body).Decode(&hostSearch)
	if err != nil {
		return HostSearch{}, err
	}
	return hostSearch, nil
}
