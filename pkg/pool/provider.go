package pool

import (
	"time"

	"github.com/chtjonas/entropyd/pkg/logging"
)

type Provider interface {
	FetchEntropy(bits int) (*Entropy, error)
}

func (p *EntropyPool) Run(interval time.Duration, targetBits, maxBits int, provider Provider) {
	backoff := make(chan struct{}, 6)
	go func() {
		for range time.Tick(2 * time.Second) {
			<-backoff
		}
	}()

	for range time.Tick(interval * time.Millisecond) {
		entropyAvail := p.GetEntropyAvail()
		writeWakeupThreshold := p.GetWriteWakeupThreshold()
		if entropyAvail < writeWakeupThreshold {
			entropyAvailable, bitsNeeded := p.GetBitsNeeded(targetBits, maxBits)
			logging.Log("fetching entropy",
				logging.LogInt("entropy_avail", entropyAvailable),
				logging.LogInt("entropy_target", targetBits),
				logging.LogInt("bits_needed", bitsNeeded),
			)
			backoff <- struct{}{}
			entropy, err := provider.FetchEntropy(bitsNeeded)
			if err != nil {
				logging.Log("failed to fetch entropy",
					logging.LogError("error", err),
				)
			} else {
				logging.Log("adding entropy to kernel pool",
					logging.LogInt("sample_size", entropy.Count),
				)
				p.AddEntropy(entropy)
			}
		}
	}
}
