package entropy

import (
	"crypto/tls"
	"net/http"
	"time"
)

type EntropyClient struct {
	serverURL string
	minBits   int
	maxBits   int
	userAgent string
	client    *http.Client
}

func NewClient(serverURL string, minBits int, maxBits int, userAgent string) *EntropyClient {
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
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    tlsconf,
	}
	cl := &http.Client{Transport: tr}

	return &EntropyClient{
		serverURL: serverURL,
		minBits:   minBits,
		maxBits:   maxBits,
		userAgent: userAgent,
		client:    cl,
	}
}
