package entropy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
)

func (client *EntropyClient) FetchEntropy(bits int) *Sample {
	bits = int(math.Ceil(float64(bits)/float64(8))) * 8
	if bits < client.minBits {
		bits = client.minBits
	}
	if bits > client.maxBits {
		bits = client.maxBits
	}
	path := fmt.Sprintf("%d", bits)
	resp, err := http.Get(client.serverURL + path)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var sample Sample
	err = json.Unmarshal(body, &sample)
	if err != nil {
		panic(err)
	}
	return &sample
}
