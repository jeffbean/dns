package dnstest

import (
	"fmt"
	"net"
	"sync"

	"github.com/miekg/dns"
)

// A Server is a DNS server listening on a system-chosen port on the
// local loopback interface, for use in end-to-end DNS tests.
type Server struct {
	Addr string

	// TODO: we can do either TCP/UDP, support either via intention
	// default to using UDP as most clients expect this by default.
	PacketConn net.PacketConn

	Config *dns.Server

	// Close blocks until all requests are finished.
	wg sync.WaitGroup
}

// NewServer starts and returns a new Server.
// The caller should call Close when finished, to shut it down.
func NewServer(handler dns.Handler) *Server {
	ts := NewUnstartedServer(handler)
	ts.Start()
	return ts
}

// NewUnstartedServer returns a new Server but doesn't start it.
//
// After changing its configuration, the caller should call Start
//
// The caller should call Close when finished, to shut it down.
func NewUnstartedServer(handler dns.Handler) *Server {
	return &Server{
		PacketConn: newLocalListener(),
		Config:     &dns.Server{Handler: handler},
	}
}

// Start starts a server from NewUnstartedServer.
func (s *Server) Start() {
	s.Addr = s.PacketConn.LocalAddr().String()
	s.goServe()
}

// Close shuts down the server and blocks until all outstanding
// requests on this server have completed.
func (s *Server) Close() {
	s.PacketConn.Close()
	s.Config.Shutdown()

	s.wg.Wait()
}

func (s *Server) goServe() {
	s.wg.Add(1)
	// Activate uses the set listener to serve from
	s.Config.PacketConn = s.PacketConn
	go func() {
		defer s.wg.Done()
		s.Config.ActivateAndServe()
	}()
}

func newLocalListener() net.PacketConn {
	l, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		panic(fmt.Sprintf("dnstest: failed to listen on a port: %v", err))
	}
	return l
}
