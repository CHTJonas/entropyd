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
	minBitsPtr := flag.Int("min", 512, "minimum amount of entropy (in bits) to call for in a single HTTP request")
	maxBitsPtr := flag.Int("max", 8192, "maximum amount of entropy (in bits) to call for in a single HTTP request")
	flag.Parse()

	cl := entropy.NewClient(*serverURLPtr, *minBitsPtr, *maxBitsPtr)
	pl := pool.OpenPool()
	defer pl.Cleardown()

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			entropyAvail := pool.GetEntropyAvail()
			writeWakeupThreshold := pool.GetWriteWakeupThreshold()
			if entropyAvail < writeWakeupThreshold {
				poolsize := pool.GetPoolsize()
				bitsNeeded := poolsize - entropyAvail

				fmt.Println("Available entropy:", entropyAvail)
				fmt.Println("Target entropy:", poolsize)
				fmt.Println("Difference:", bitsNeeded)

				sample, err := cl.FetchEntropy(bitsNeeded)
				if err != nil {
					fmt.Println(err)
				} else {
					if sample.Validate() {
						pl.AddEntropy(sample)
					} else {
						fmt.Println("Entropy sample failed validation")
					}
				}

				fmt.Println()
			}
		}
	}
}
