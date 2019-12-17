package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
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

type randPoolInfo struct {
	entropyCount int
	bufSize      int
	buf          []byte
}

const entropyServerURL = "https://entropy.malc.org.uk/entropy/"
const rndAddEntropy = 0x40085203
const minBits = 512
const maxBits = 8192

func fetchEntropy(bits int) *Sample {
	bits = int(math.Ceil(float64(bits)/float64(8))) * 8
	if bits < minBits {
		bits = minBits
	}
	if bits > maxBits {
		bits = maxBits
	}
	path := fmt.Sprintf("%d", bits)
	resp, err := http.Get(entropyServerURL + path)
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
	arg := unsafe.Pointer(&randPoolInfo{
		entropyCount: sample.getBits(),
		bufSize:      sample.getSize(),
		buf:          sample.getData(),
	})
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(rndAddEntropy), uintptr(arg))
	if ep != 0 {
		err := syscall.Errno(ep)
		panic(err)
	}
}

func getEntropyAvail() int {
	path := "/proc/sys/kernel/random/entropy_avail"
	return readIntFromFile(path)
}

func getWriteWakeupThreshold() int {
	path := "/proc/sys/kernel/random/write_wakeup_threshold"
	return readIntFromFile(path)
}

func getPoolsize() int {
	path := "/proc/sys/kernel/random/poolsize"
	return readIntFromFile(path)
}

func readIntFromFile(path string) int {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	i, err := strconv.Atoi(scanner.Text())

	if err != nil {
		panic(err)
	}
	return i
}

func main() {
	fd, err := syscall.Open("/dev/random", syscall.O_RDWR, 666)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(fd)
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			entropyAvail := getEntropyAvail()
			writeWakeupThreshold := getWriteWakeupThreshold()
			if entropyAvail < writeWakeupThreshold {
				poolsize := getPoolsize()
				bitsNeeded := poolsize - entropyAvail
				sample := fetchEntropy(bitsNeeded)
				if sample.validate() {
					addEntropy(sample, fd)
				}
			}
		}
	}
}
