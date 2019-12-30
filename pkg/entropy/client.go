package entropy

import (
	"net/http"
	"time"
)

type EntropyClient struct {
	serverURL string
	minBits   int
	maxBits   int
	client    *http.Client
}

func NewClient(serverURL string, minBits int, maxBits int) *EntropyClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	cl := &http.Client{Transport: tr}

	return &EntropyClient{
		serverURL: serverURL,
		minBits:   minBits,
		maxBits:   maxBits,
		client:    cl,
	}
}
