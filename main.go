package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"syscall"
	"unsafe"
)

type Sample struct {
	Entropy Entropy `json:"entropy"`
}

func (sample *Sample) validate() bool {
	return sample.Entropy.validate()
}

func (sample *Sample) getSize() int {
	return sample.Entropy.getSize()
}

func (sample *Sample) getBits() int {
	return sample.Entropy.getBits()
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

func (entropy *Entropy) getSize() int {
	return entropy.Length
}

func (entropy *Entropy) getBits() int {
	return entropy.Bits
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

type rand_pool_info struct {
	entropy_count int
	buf_size      int
	buf           []byte
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

func addEntropy(sample *Sample, fd int) {
	RNDADDENTROPY := 0x40085203
	arg := unsafe.Pointer(&rand_pool_info{
		entropy_count: sample.getBits(),
		buf_size:      sample.getSize(),
		buf:           sample.getData(),
	})
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(RNDADDENTROPY), uintptr(arg))
	if ep != 0 {
		err := syscall.Errno(ep)
		panic(err)
	}
}

func main() {
	fd, err := syscall.Open("/dev/random", syscall.O_RDWR, 666)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)
	sample := fetchEntropy(4096)
	if sample.validate() {
		addEntropy(sample, fd)
	}
}
