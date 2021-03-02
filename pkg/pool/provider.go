package pool

import (
	"time"

	"github.com/chtjonas/entropyd/pkg/logging"
)

type Provider interface {
	FetchEntropy(bits int) (*Entropy, error)
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

	return nil
}
