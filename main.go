package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
)

func resolve(name string, nameserver string) (net.IP, error) {
	for {
		reply := dnsQuery(name, nameserver)
		if reply == nil {
			return nil, fmt.Errorf("resolve server %s could not be reached", nameserver)
		}
		ip := getAnswer(reply)
		if ip != nil {
			// Best case: we get an answer to our query and we're done
			return ip, nil
		}
		nsIP := getAdditional(reply)
		if nsIP != nil {
			// Second best: we get a "glue record" in the additional section with the *IP address* of another nameserver to query
			nameserver = nsIP.String()
		} else {
			domain := getAuthority(reply)
			if domain != "" {
				// Third best: we get the *domain name* of another nameserver in the NS / authority / extra section to query, which we can look up the IP for
				ip, err := resolve(domain, nameserver)
				if err != nil {
					return nil, err
				}
				nameserver = ip.String()
			} else {
				// If there's no A record after all tries, return an error
				return nil, fmt.Errorf("failed to resolve %s ", name)
			}
		}
	}
}

func getAnswer(reply *dns.Msg) net.IP {
	for _, record := range reply.Answer {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

func getAdditional(reply *dns.Msg) net.IP {
	for _, record := range reply.Extra {
		if record.Header().Rrtype == dns.TypeA {
			fmt.Println("  ", record)
			return record.(*dns.A).A
		}
	}
	return nil
}

func getAuthority(reply *dns.Msg) string {
	for _, record := range reply.Ns {
		if record.Header().Rrtype == dns.TypeNS {
			fmt.Println("  ", record)
			return record.(*dns.NS).Ns
		}
		if record.Header().Rrtype == 6 {
			fmt.Println("  ", record)
			return record.(*dns.SOA).Ns
		}
	}
	return ""
}

func dnsQuery(name string, server string) *dns.Msg {
	fmt.Printf("dig -r @%s %s\n", server, name)
	msg := new(dns.Msg)
	msg.SetQuestion(name, dns.TypeA)
	c := new(dns.Client)
	reply, _, err := c.Exchange(msg, server+":53")
	if (err != nil) {
		fmt.Printf("Error querying %s: %v\n", server, err)
	}
	return reply
}


func main() {
	// https://jvns.ca/blog/2022/02/01/a-dns-resolver-in-80-lines-of-go/
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <domain> [<dns server>]")
		os.Exit(1)
	}

	name := dns.Fqdn(os.Args[1])   // Fully qualify the domain
	nameserver := ""               // Default DNS server

	if len(os.Args) > 2 {
		fmt.Println("Using resolver:", os.Args[2])
		nameserver = os.Args[2]
	} else {
		config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
    	if err != nil {
        	fmt.Println("Error loading default resolver config:", err)
        	nameserver = "8.8.8.8"
    	} else {
			fmt.Println("Using system resolver:", config.Servers[0])
			nameserver = config.Servers[0]
		}
	}

	if !strings.HasSuffix(name, ".") {
		name = name + "."
	}
	ip, err := resolve(name, nameserver)
	if (err != nil) {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", ip)
	}
}
