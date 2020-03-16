package types

import (
	"testing"
)

func TestMarshalAndUnmarshalData(t *testing.T) {
	// TODO:
	// addr := GenerateEthAddress()
	// hash := ethcmn.BigToHash(big.NewInt(2))

	// e := encodableTxData{
	// 	AccountNonce: 2,
	// 	Price:        utils.MarshalBigInt(big.NewInt(3)),
	// 	GasLimit:     1,
	// 	Recipient:    &addr,
	// 	Amount:       utils.MarshalBigInt(big.NewInt(4)),
	// 	Payload:      []byte("test"),

	// 	V: utils.MarshalBigInt(big.NewInt(5)),
	// 	R: utils.MarshalBigInt(big.NewInt(6)),
	// 	S: utils.MarshalBigInt(big.NewInt(7)),

	// 	Hash: &hash,
	// }
	// str, err := marshalAmino(e)
	// require.NoError(t, err)

	// e2 := new(encodableTxData)

	// err = unmarshalAmino(e2, str)
	// require.NoError(t, err)
	// require.Equal(t, e, *e2)
}
