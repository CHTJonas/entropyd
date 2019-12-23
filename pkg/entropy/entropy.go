package entropy

import "encoding/base64"

type Entropy struct {
	Data    string `json:"data_b64"`
	Length  int    `json:"data_len"`
	Bits    int    `json:"entropy_bits"`
	Magic   int    `json:"random_magic"`
	Source  string `json:"source"`
	Version int    `json:"version"`
}

func (entropy *Entropy) GetSize() int {
	return entropy.Length
}

func (entropy *Entropy) GetBits() int {
	return entropy.Bits
}

func (entropy *Entropy) GetData() []byte {
	decoded, err := base64.StdEncoding.DecodeString(entropy.Data)
	if err != nil {
		panic(err)
	}
	return decoded
}

func (entropy *Entropy) Validate() bool {
	data := entropy.GetData()
	if len(data) != entropy.Length {
		return false
	}
	if 8*len(data) < entropy.Bits {
		return false
	}
	return true
}
