package types

import "math/big"

// MarshalBigInt marshals big int into text string for consistent encoding
func MarshalBigInt(i *big.Int) (string, error) {
	bz, err := i.MarshalText()
	if err != nil {
		return "", err
	}
	return string(bz), nil
}

// UnmarshalBigInt unmarshals string from *big.Int
func UnmarshalBigInt(s string) (*big.Int, error) {
	ret := new(big.Int)
	err := ret.UnmarshalText([]byte(s))
	if err != nil {
		return nil, err
	}
	return ret, nil
}
