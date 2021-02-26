package malc

import (
	"encoding/base64"
	"errors"
)

type Sample struct {
	Data    string `json:"data_b64"`
	Length  int    `json:"data_len"`
	Bits    int    `json:"entropy_bits"`
	Magic   int    `json:"random_magic"`
	Source  string `json:"source"`
	Version int    `json:"version"`
}

func (s *Sample) GetData() []byte {
	decoded, err := base64.StdEncoding.DecodeString(s.Data)
	if err != nil {
		panic(err)
	}
	return decoded
}

func (s *Sample) Validate(length int) error {
	if length != s.Length {
		return errors.New("Bad data length")
	}
	if 8*length < s.Bits {
		return errors.New("Server claims impossibly-good entropy quality")
	}
	return nil
}
