package pool

import "errors"

type Entropy struct {
	Data  []byte
	Count int
}

var ErrImpossiblyHighQuality = errors.New("Claimed entropy quality is impossibly high")

func (e *Entropy) Validate() error {
	if e.Count > 8*len(e.Data) {
		return ErrImpossiblyHighQuality
	}
	return nil
}
