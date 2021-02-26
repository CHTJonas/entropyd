package malc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/chtjonas/entropyd/pkg/pool"
)

func (client *EntropyClient) FetchEntropy(bits int) (*pool.Entropy, error) {
	bits = roundToOctet(bits)
	bits = clamp(bits, client.minBits, client.maxBits)
	body, err := client.requestFromServer(bits)
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

func (client *EntropyClient) requestFromServer(bits int) ([]byte, error) {
	url := fmt.Sprintf("%s%d", client.serverURL, bits)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
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
