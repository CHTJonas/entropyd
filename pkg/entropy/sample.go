package entropy

type Sample struct {
	Entropy Entropy `json:"entropy"`
}

func (sample *Sample) GetSize() int {
	return sample.Entropy.GetSize()
}

func (sample *Sample) GetBits() int {
	return sample.Entropy.GetBits()
}

func (sample *Sample) GetData() []byte {
	return sample.Entropy.GetData()
}

func (sample *Sample) Validate() bool {
	return sample.Entropy.Validate()
}
