package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
)

func main() {
	var msg dns.Msg
	fqdn := dns.Fqdn("stacktitan.com")
	msg.SetQuestion(fqdn, dns.TypeA)

	resp, err := dns.Exchange(&msg, "8.8.8.8:53")
	if err != nil {
		log.Fatal(err)
	}
	if len(resp.Answer) < 1 {
		fmt.Println("No answer")
		return
	}

	for _, answer := range resp.Answer {
		if a, ok := answer.(*dns.A); ok {
			fmt.Println(a.A)
		}
	}
}
