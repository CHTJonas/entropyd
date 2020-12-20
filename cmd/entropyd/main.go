package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/chtjonas/entropyd/pkg/entropy"
	"github.com/chtjonas/entropyd/pkg/pool"
)

// Limit of ioctl requests is 1024 bytes, including header.
const maxDataBytes = 1016

// Must be no higher than that used by the server.
const maxReqBits = maxDataBytes * 8

func main() {
	checkOS()

	ipv4Ptr := flag.Bool("4", false, "force the use of IPv4 for the HTTP connection")
	ipv6Ptr := flag.Bool("6", false, "force the use of IPv6 for the HTTP connection")
	serverURLPtr := flag.String("url", "https://entropy.malc.org.uk/entropy/", "URL of the remote entropy server")
	minBitsPtr := flag.Int("min", 64, "minimum amount of entropy (in bits) in a HTTP request")
	maxBitsPtr := flag.Int("max", maxReqBits, "maximum amount of entropy (in bits) in a HTTP request")
	targetBitsPtr := flag.Int("target", 3072, "target amount of entropy (in bits) to store in the kernel entropy pool")
	pollIntervalPtr := flag.Int("poll", 200, "interval (in milliseconds) at which to poll the kernel entropy pool")
	doDryRunPtr := flag.Bool("dry-run", false, "makes a request for 512 bits of entropy but writes to stdout instead of the kernel entropy pool")
	flag.Parse()

	ver := getVer().getString()
	ua := "entropyd/" + ver + " (+https://github.com/CHTJonas/entropyd)"
	ipv := ""
	if *ipv4Ptr {
		ipv = "tcp4"
	}
	if *ipv6Ptr {
		ipv = "tcp6"
	}
	cl := entropy.NewClient(*serverURLPtr, *minBitsPtr, *maxBitsPtr, ua, ipv)
	pl := pool.OpenPool()
	defer pl.Cleardown()

	if *doDryRunPtr {
		sample, err := cl.FetchEntropy(16)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data := sample.GetData()
		fmt.Printf("Entropy: %s", data)
		os.Exit(0)
	}

	interval := time.Duration(*pollIntervalPtr)
	backoff := make(chan interface{}, 6)

	go func() {
		for range time.Tick(2 * time.Second) {
			<-backoff
		}
	}()

	for range time.Tick(interval * time.Millisecond) {
		entropyAvail := pl.GetEntropyAvail()
		writeWakeupThreshold := pl.GetWriteWakeupThreshold()
		if entropyAvail < writeWakeupThreshold {
			entropyAvailable, bitsNeeded := pl.GetBitsNeeded(*targetBitsPtr, *maxBitsPtr)
			fmt.Printf("Entropy available: %d. Entropy target: %d. Entropy needed: %d.\n", entropyAvailable, *targetBitsPtr, bitsNeeded)
			backoff <- nil
			sample, err := cl.FetchEntropy(bitsNeeded)
			if err != nil {
				fmt.Println(err)
			} else {
				err := sample.Validate()
				if err != nil {
					fmt.Println(err)
				} else {
					pl.AddEntropy(sample)
				}
			}
		}
	}
}

func checkOS() {
	if runtime.GOOS != "linux" {
		fmt.Println("entropyd can only run on Linux")
		os.Exit(1)
	}
}
