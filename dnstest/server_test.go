package dnstest

import (
	"testing"

	"github.com/miekg/dns"
)

func TestNewServer(t *testing.T) {
	s := NewServer(dns.DefaultServeMux)
	defer s.Close()
}
