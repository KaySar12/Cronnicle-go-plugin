package dnsclient

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
)

type DNSClient struct{}

func LoadDnsServer() (*dns.ClientConfig, error) {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to load system DNS configuration: %w", err)
	}
	config.Port = "53"
	return config, nil
}

func Lookup(config *dns.ClientConfig, lookupType uint16, zone string, timeout int) (*dns.Msg, error) {
	dnsServer := net.JoinHostPort(config.Servers[0], config.Port)
	msg := new(dns.Msg)
	msg.SetQuestion(zone, lookupType) // Query for NS records
	msg.RecursionDesired = true
	client := new(dns.Client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	response, _, err := client.ExchangeContext(ctx, msg, dnsServer)
	if err != nil {
		return nil, fmt.Errorf("failed to query DNS for NS records: %w", err)
	}
	return response, nil
}

func FormatZone(zone string) string {
	if zone[len(zone)-1] != '.' {
		zone += "."
	}
	return zone
}
