package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/chtjonas/entropy-client/pkg/entropy"
	"github.com/chtjonas/entropy-client/pkg/pool"
)

// Limit of ioctl requests is 1024 bytes, including header.
const maxDataBytes = 1016

// Must be no higher than that used by the server.
const maxReqBits = maxDataBytes * 8

func main() {
	serverURLPtr := flag.String("url", "https://entropy.malc.org.uk/entropy/", "URL of the remote entropy server")
	minBitsPtr := flag.Int("min", 64, "minimum amount of entropy (in bits) in a HTTP request")
	maxBitsPtr := flag.Int("max", maxReqBits, "maximum amount of entropy (in bits) in a HTTP request")
	targetBitsPtr := flag.Int("target", 3072, "target amount of entropy (in bits) to store in the kernel entropy pool")
	pollIntervalPtr := flag.Int("poll", 200, "interval (in milliseconds) at which to poll the kernel entropy pool")
	doDryRunPtr := flag.Bool("dry-run", false, "makes a request for 512 bits of entropy but writes to stdout instead of the kernel entropy pool")
	flag.Parse()

	ver := getVer().getString()
	ua := "entropy-client/" + ver + " (+https://github.com/CHTJonas/entropy-client)"
	cl := entropy.NewClient(*serverURLPtr, *minBitsPtr, *maxBitsPtr, ua)
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
	ticker := time.NewTicker(interval * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			entropyAvail := pl.GetEntropyAvail()
			writeWakeupThreshold := pl.GetWriteWakeupThreshold()
			if entropyAvail < writeWakeupThreshold {
				bitsNeeded := pl.GetBitsNeeded(*targetBitsPtr, *maxBitsPtr)
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
}
