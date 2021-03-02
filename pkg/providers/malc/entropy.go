package malc

import (
	"encoding/base64"
	"errors"
)

var ErrBadLength = errors.New("Bad data length")
var ErrImpossibleQuality = errors.New("Server claims impossibly-good entropy quality")

type Sample struct {
	Data    string `json:"data_b64"`
	Length  int    `json:"data_len"`
	Bits    int    `json:"entropy_bits"`
	Magic   int    `json:"random_magic"`
	Source  string `json:"source"`
	Version int    `json:"version"`
}

func (s *Sample) GetData() ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(s.Data)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func (s *Sample) Validate(length int) error {
	if length != s.Length {
		return ErrBadLength
	}
	if 8*length < s.Bits {
		return ErrImpossibleQuality
	}
	return nil
}
