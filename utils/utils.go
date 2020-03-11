package utils

import (
	"encoding/hex"
	"errors"
	"strings"
)

// HexToBytes turns a 0x prefixed hex string into a byte slice
func HexToBytes(in string) ([]byte, error) {
	if len(in) < 2 {
		return nil, errors.New("invalid string")
	}

	if strings.Compare(in[:2], "0x") != 0 {
		return nil, errors.New("could not byteify non 0x prefixed string")
	}
	// Ensure we have an even length, otherwise hex.DecodeString will fail and return zero hash
	if len(in)%2 != 0 {
		return nil, errors.New("cannot decode a odd length string")
	}
	in = in[2:]
	out, err := hex.DecodeString(in)
	return out, err
}
