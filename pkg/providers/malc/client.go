package malc

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

const ServerURL = "https://entropy.malc.org.uk/entropy/"

type EntropyClient struct {
	minBits   int
	maxBits   int
	userAgent string
	client    *http.Client
}

func NewEntropyClient(minBits int, maxBits int, userAgent, ipVersion string) *EntropyClient {
	tlsconf := &tls.Config{
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
		Timeout:   30 * time.Second,
		KeepAlive: 10 * time.Second,
	}
	dialCtx := func(ctx context.Context, network, addr string) (net.Conn, error) {
		if ipVersion != "" {
			network = ipVersion
		}
		return dialer.DialContext(ctx, network, addr)
	}
	tr := &http.Transport{
		MaxIdleConns:        10,
		IdleConnTimeout:     10 * time.Minute,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
		DisableCompression:  true,
		TLSClientConfig:     tlsconf,
		DialContext:         dialCtx,
	}
	cl := &http.Client{Transport: tr}

	return &EntropyClient{
		minBits:   minBits,
		maxBits:   maxBits,
		userAgent: userAgent,
		client:    cl,
	}
}
