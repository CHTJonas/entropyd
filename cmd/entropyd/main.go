package main

import (
	"fmt"
	"time"

	"github.com/chtjonas/entropy-client/pkg/entropy"
	"github.com/chtjonas/entropy-client/pkg/pool"
)

const serverURL = "https://entropy.malc.org.uk/entropy/"
const minBits = 512
const maxBits = 8192

func main() {
	cl := entropy.NewClient(serverURL, minBits, maxBits)
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

				fmt.Println(entropyAvail)
				fmt.Println(writeWakeupThreshold)
				fmt.Println(bitsNeeded)

				sample := cl.FetchEntropy(bitsNeeded)
				if sample.Validate() {
					pl.AddEntropy(sample)
				}

				fmt.Println("")
			}
		}
	}
}
