package types

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common/bitutil"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Code is account Code.
type Code struct {
	// Compressed bytes using the sparse bitset representation algorithm. Hex
	// formatted for readability in genesis state.
	CompressedBytes hexutil.Bytes `json:"compressed_bytes"`
	// Size defines the original bytes size. It it used as the target for decompressing
	// the code.
	Size int `json:"size"`
}

// NewCode creates a new Code instance from uncompressed bytes.
func NewCode(bytes []byte) Code {
	return Code{
		CompressedBytes: bitutil.CompressBytes(bytes),
		Size:            len(bytes),
	}
}

// Bytes returns the decompressed bytes from the Code. It uses the original byte
// length as the target size. If the compression fails, the function will panic.
func (c Code) Bytes() []byte {
	bz, err := bitutil.DecompressBytes(c.CompressedBytes, c.Size)
	if err != nil {
		panic(err)
	}

	return bz
}

// String returns the compressed bytes in hex format.
func (c Code) String() string {
	return c.CompressedBytes.String()
}

// IsEmpty returns true if the code is uninitialized or contains empty values.
func (c Code) IsEmpty() bool {
	return c.Size == 0 && len(c.CompressedBytes) == 0
}

// Validate checks that the bytes are not empty and that the original size is greater
// or equal than the compressed bytes length.
func (c Code) Validate() error {
	if c.IsEmpty() {
		return errors.New("code bytes cannot be empty")
	}
	if c.Size < len(c.CompressedBytes) {
		return fmt.Errorf("original code size should be ≥ compressed bytes size")
	}
	return nil
}
