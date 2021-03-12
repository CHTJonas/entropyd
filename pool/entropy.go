package pool

import "errors"

type Entropy struct {
	Data  []byte
	Count int
}

var ErrDataTooLarge = errors.New("entropy payload is too large")
var ErrImpossiblyHighQuality = errors.New("claimed entropy quality is impossibly high")

func (e *Entropy) Validate() error {
	length := len(e.Data)
	if length > MaxDataBytes {
		return ErrDataTooLarge
	}
	if e.Count > 8*length {
		return ErrImpossiblyHighQuality
	}
	return nil
}
