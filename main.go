package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall"
)

type Sample struct {
	Entropy Entropy `json:"entropy"`
}

func (sample *Sample) validate() bool {
	return sample.Entropy.validate()
}

func (sample *Sample) getData() []byte {
	return sample.Entropy.getData()
}

type Entropy struct {
	Data    string `json:"data_b64"`
	Length  int    `json:"data_len"`
	Bits    int    `json:"entropy_bits"`
	Magic   int    `json:"random_magic"`
	Source  string `json:"source"`
	Version int    `json:"version"`
}

func (entropy *Entropy) getData() []byte {
	decoded, err := base64.StdEncoding.DecodeString(entropy.Data)
	if err != nil {
		panic(err)
	}
	return decoded
}

func (entropy *Entropy) validate() bool {
	data := entropy.getData()
	if len(data) != entropy.Length {
		return false
	}
	if 8*len(data) < entropy.Bits {
		return false
	}
	return true
}

func fetchEntropy(bits uint) *Sample {
	entropyServerURL := "https://entropy.malc.org.uk/entropy/"
	bitPath := fmt.Sprintf("%d", bits)
	resp, err := http.Get(entropyServerURL + bitPath)
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

func fillPool(data []byte, fd int) {
	fmt.Println(fd)
	_, err := syscall.Write(fd, data)
	if err != nil {
		panic(err)
	}
}

func main() {
	// fd, err := syscall.Open("/dev/random", syscall.O_RDWR, 666)
	fd, err := syscall.Open("/tmp/random", syscall.O_RDWR, 666)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)
	sample := fetchEntropy(4096)
	if sample.validate() {
		fillPool(sample.getData(), fd)
	}
}
