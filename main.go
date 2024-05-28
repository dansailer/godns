package main

import (
	"fmt"
	"net"
	"os"

	"github.com/miekg/dns"
)

func main() {
	if len(os.Args) < 2 || len(os.Args) > 4 {
		fmt.Println("Usage: go run main.go <domain> [<record type>] [<dns server>]")
		os.Exit(1)
	}

	domain := dns.Fqdn(os.Args[1]) // Fully qualify the domain
	recordType := "A"              // Default record type
	dnsServer := ""                // Default DNS server

	if len(os.Args) >= 3 {
		recordType = os.Args[2]
	}

	if len(os.Args) == 4 {
		dnsServer = os.Args[3]
	} else {
		// Use net.LookupIP to trigger the use of the system's default DNS server
		_, err := net.LookupIP("example.com")
		if err != nil {
			fmt.Println("Failed to use system default DNS server, using 8.8.8.8:53")
			dnsServer = "8.8.8.8:53"
		} else {
			// Use a default DNS server if not specified
			dnsServer = "8.8.8.8:53"
		}
	}

	c := new(dns.Client)
	m := new(dns.Msg)

	qtype, ok := dns.StringToType[recordType]
	if !ok {
		fmt.Println("Invalid record type:", recordType)
		os.Exit(1)
	}

	m.SetQuestion(domain, qtype)

	r, _, err := c.Exchange(m, dnsServer)
	if err != nil {
		fmt.Println("DNS query failed:", err)
		os.Exit(1)
	}

	if len(r.Answer) == 0 {
		fmt.Println("No records found.")
		return
	}

	for _, ans := range r.Answer {
		fmt.Println(ans)
	}
}
