package malc

import "github.com/chtjonas/entropyd/pkg/pool"

type Wrapper struct {
	Sample Sample `json:"entropy"`
}

func (w *Wrapper) ToEntropy() (*pool.Entropy, error) {
	e := &pool.Entropy{
		Data:  w.Sample.GetData(),
		Count: w.Sample.Bits,
	}
	length := len(e.Data)
	if err := w.Sample.Validate(length); err != nil {
		return nil, err
	}
	return e, nil
}
