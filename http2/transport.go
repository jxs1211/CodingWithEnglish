package http2

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
)

var NextProtoTLS string

type Transport struct {
	Fallback http.RoundTripper
}

func (t *Transport) RoundTripper(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme != "https" {
		if t.Fallback == nil {
			return nil, errors.New("http2: unsupport protocol")
		}
		return t.Fallback.RoundTrip(req)
	}
	host, port, err := net.SplitHostPort(req.URL.Host)
	if err != nil {
		host = req.URL.Host
		port = "443"
	}
	cfg := &tls.Config{
		NextProtos: []string{NextProtoTLS},
	}
	tConn, err := tls.Dial("tcp", fmt.Sprintf("%v:%v", host, port), cfg)
	if err != nil {
		return nil, err
	}
	err = tConn.Handshake()
	if err != nil {
		return nil, err
	}
	state := tConn.ConnectionState()
	if state.NegotiatedProtocol != NextProtoTLS {
		return nil, fmt.Errorf("bad protocol")
	}
	if !state.NegotiatedProtocolIsMutual {
		return nil, errors.New("could not negotiate protocol mutually")
	}
	return nil, errors.New("TODO")
}
