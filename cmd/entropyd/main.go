package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/chtjonas/entropy-client/pkg/entropy"
	"github.com/chtjonas/entropy-client/pkg/pool"
)

func main() {
	serverURLPtr := flag.String("url", "https://entropy.malc.org.uk/entropy/", "URL of the remote entropy server")
	minBitsPtr := flag.Int("min", 512, "minimum amount of entropy (in bits) in a HTTP request")
	maxBitsPtr := flag.Int("max", 8192, "maximum amount of entropy (in bits) in a HTTP request")
	pollIntervalPtr := flag.Int("poll", 200, "interval (in milliseconds) at which to poll the kernel entropy pool")
	flag.Parse()

	cl := entropy.NewClient(*serverURLPtr, *minBitsPtr, *maxBitsPtr)
	pl := pool.OpenPool()
	defer pl.Cleardown()

	interval := time.Duration(*pollIntervalPtr)
	ticker := time.NewTicker(interval * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			entropyAvail := pool.GetEntropyAvail()
			writeWakeupThreshold := pool.GetWriteWakeupThreshold()
			if entropyAvail < writeWakeupThreshold {
				poolSize := pool.GetPoolSize()
				bitsNeeded := poolSize - entropyAvail
				fmt.Printf("Entropy available: %d. Entropy target: %d. Entropy delta: %d.\n", entropyAvail, poolSize, bitsNeeded)
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
