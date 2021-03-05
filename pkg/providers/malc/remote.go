package malc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/chtjonas/entropyd/pkg/pool"
)

var seedOnce sync.Once

func (client *EntropyClient) FetchEntropy(bits int) (*pool.Entropy, error) {
	bits = roundToOctet(bits)
	bits = clamp(bits, client.minBits, client.maxBits)
	body, err := client.requestFromServers(bits)
	if err != nil {
		return nil, err
	}
	w := new(Wrapper)
	err = json.Unmarshal(body, w)
	if err != nil {
		return nil, err
	}
	return w.ToEntropy()
}

func (client *EntropyClient) requestFromServers(bits int) ([]byte, error) {
	ips, err := net.LookupIP(ServerName)
	if err != nil {
		return nil, err
	}
	seedOnce.Do(func() { rand.Seed(time.Now().UnixNano()) })
	rand.Shuffle(len(ips), func(i, j int) { ips[i], ips[j] = ips[j], ips[i] })
	for _, ip := range ips {
		var data []byte
		data, err = client.makeRequest(bits, ip)
		if err != nil {
			continue
		}
		return data, nil
	}
	return nil, err
}

func (client *EntropyClient) makeRequest(bits int, ip net.IP) ([]byte, error) {
	ipStr := ip.String()
	if !isIPv4(ip) {
		ipStr = "[" + ipStr + "]"
	}
	url := fmt.Sprintf("https://%s/entropy/%d", ipStr, bits)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Host = ServerName
	req.Header.Set("Cache-Control", "no-store, max-age=0")
	if client.userAgent != "" {
		req.Header.Set("User-Agent", client.userAgent)
	}
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func isIPv4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}
