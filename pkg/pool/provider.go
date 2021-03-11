package pool

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Provider interface {
	FetchEntropy(bits int) (*Entropy, error)
}

func (p *EntropyPool) SetLogger(logger *log.Logger) {
	p.logger = logger
}

func (p *EntropyPool) log(msg string, tuples map[string]interface{}) {
	if p.logger == nil {
		return
	}
	b := new(strings.Builder)
	fmt.Fprintf(b, "msg=%s", msg)
	for k, v := range tuples {
		fmt.Fprintf(b, ", %s=%v", k, v)
	}
	p.logger.Println(b.String())
}

func (p *EntropyPool) Run(interval time.Duration, targetBits, maxBits int, provider Provider) error {
	backoff := make(chan struct{}, 6)
	go func() {
		for range time.Tick(2 * time.Second) {
			<-backoff
		}
	}()

	for range time.Tick(interval * time.Millisecond) {
		entropyAvail, err := p.GetEntropyAvail()
		if err != nil {
			return err
		}
		writeWakeupThreshold, err := p.GetWriteWakeupThreshold()
		if err != nil {
			return err
		}
		if entropyAvail < writeWakeupThreshold {
			entropyAvailable, bitsNeeded, err := p.GetBitsNeeded(targetBits, maxBits)
			if err != nil {
				return err
			}
			p.log("fetching entropy from provider", map[string]interface{}{
				"entropy_avail":  entropyAvailable,
				"entropy_target": targetBits,
				"bits_needed":    bitsNeeded,
			})
			backoff <- struct{}{}
			entropy, err := provider.FetchEntropy(bitsNeeded)
			if err != nil {
				p.log("failed to fetch entropy", map[string]interface{}{
					"error": err,
				})
			} else {
				p.log("adding entropy to kernel pool", map[string]interface{}{
					"sample_size": entropy.Count,
				})
				err = p.AddEntropy(entropy)
				if err != nil {
					p.log("failed to add entropy", map[string]interface{}{
						"error": err,
					})
				}
			}
		}
	}
	return nil
}
