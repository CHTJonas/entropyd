package malc

import "github.com/chtjonas/entropyd/pkg/pool"

type Wrapper struct {
	Sample Sample `json:"entropy"`
}

func (w *Wrapper) ToEntropy() (*pool.Entropy, error) {
	data, err := w.Sample.GetData()
	if err != nil {
		return nil, err
	}
	count := w.Sample.Bits
	e := &pool.Entropy{
		Data:  data,
		Count: count,
	}
	length := len(e.Data)
	if err := w.Sample.Validate(length); err != nil {
		return nil, err
	}
	return e, nil
}
