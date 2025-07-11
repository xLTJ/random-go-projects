package dnsStuff

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
	"strings"
)

type Result struct {
	IPAddress string
	Hostname  string
}

func Lookup(fqdn string, conn net.Conn) ([]Result, error) {
	var results []Result
	var currentFqdn = dns.Fqdn(fqdn)
	for {
		cname, err := LookupCNAME(fqdn, conn)
		if err == nil && len(cname) > 0 {
			currentFqdn = cname[0]
			continue
		}

		ips, err := LookupA(currentFqdn, conn)
		if err != nil {
			//return nil, fmt.Errorf("error shit: %v", err) // uncomment to see errors
			if strings.Contains(err.Error(), "i/o timeout") {
				continue // try again
			}
			break // other error is probably that there was no result and the subdomain is invalid
		}

		for _, ip := range ips {
			results = append(results, Result{
				IPAddress: ip,
				Hostname:  fqdn + "---" + conn.RemoteAddr().String(),
			})
		}
		break
	}

	return results, nil
}

func LookupA(fqdn string, conn net.Conn) ([]string, error) {
	var msg dns.Msg
	var ips []string
	var client dns.Client
	msg.SetQuestion(fqdn, dns.TypeA)

	resp, _, err := client.ExchangeWithConn(&msg, &dns.Conn{Conn: conn})
	if err != nil {
		return ips, fmt.Errorf("error retrieving A records: %v", err)
	}
	if len(resp.Answer) == 0 {
		return ips, fmt.Errorf("no answer")
	}

	for _, answer := range resp.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

func LookupCNAME(fqdn string, conn net.Conn) ([]string, error) {
	var msg dns.Msg
	var fqdns []string
	var client dns.Client
	msg.SetQuestion(fqdn, dns.TypeCNAME)

	resp, _, err := client.ExchangeWithConn(&msg, &dns.Conn{Conn: conn})
	if err != nil {
		return fqdns, fmt.Errorf("error retrieving CNAME records: %v", err)
	}
	if len(resp.Answer) == 0 {
		return fqdns, fmt.Errorf("no answer")
	}

	for _, answer := range resp.Answer {
		if cname, ok := answer.(*dns.CNAME); ok {
			fqdns = append(fqdns, cname.Target)
		}
	}
	return fqdns, nil
}
