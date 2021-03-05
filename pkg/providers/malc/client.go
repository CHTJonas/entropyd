package malc

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const ServerName = "entropy.malc.org.uk"

type EntropyClient struct {
	minBits   int
	maxBits   int
	userAgent string
	client    *http.Client
}

func NewEntropyClient(minBits int, maxBits int, userAgent, ipVersion string) *EntropyClient {
	tlsconf := &tls.Config{
		ServerName:               ServerName,
		MinVersion:               tls.VersionTLS12,
		PreferServerCipherSuites: false,
		CurvePreferences:         []tls.CurveID{tls.X25519},
		CipherSuites: []uint16{
			tls.TLS_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
		},
	}
	dialer := &net.Dialer{
		Timeout:   5 * time.Second,
		KeepAlive: -1,
	}
	dialCtx := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if ipVersion != "" {
			network = ipVersion
		}
		return dialer.DialContext(ctx, network, addr)
	}
	tr := &http.Transport{
		TLSHandshakeTimeout: 3 * time.Second,
		DisableKeepAlives:   true,
		DisableCompression:  true,
		TLSClientConfig:     tlsconf,
		DialContext:         dialCtx,
	}
	cl := &http.Client{
		Transport: tr,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	return &EntropyClient{
		minBits:   minBits,
		maxBits:   maxBits,
		userAgent: userAgent,
		client:    cl,
	}
}
