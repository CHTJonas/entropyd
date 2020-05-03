package entropy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

func (client *EntropyClient) FetchEntropy(bits int) (*Sample, error) {
	bits = int(math.Ceil(float64(bits)/float64(8))) * 8
	if bits < client.minBits {
		bits = client.minBits
	}
	if bits > client.maxBits {
		bits = client.maxBits
	}
	body, err := client.RequestFromServer(bits)
	if err != nil {
		return nil, err
	}
	var sample Sample
	err = json.Unmarshal(body, &sample)
	if err != nil {
		return nil, err
	}
	return &sample, nil
}

func (client *EntropyClient) RequestFromServer(bits int) ([]byte, error) {
	httpclient := client.client
	path := fmt.Sprintf("%d", bits)
	url := client.serverURL + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {

	}
	if client.userAgent != "" {
		req.Header.Set("User-Agent", client.userAgent)
	}
	resp, err := httpclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
