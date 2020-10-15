package types

import (
	"bytes"
	"errors"
	"fmt"
	//"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
)

// Order contains original signed order and status fields.
type Order struct {
	Order                  *SafeSignedOrder `json:"order"`
	TradePairHash          common.Hash      `json:"pairHash"`
	FilledTakerAssetAmount BigNum           `json:"filledTakerAssetAmount"`
	Status                 OrderStatus      `json:"status"`
}

// Order contains original signed order and status fields.
type DerivativeOrder struct {
	Order            *SafeSignedOrder `json:"order"`
	DerivativeMarket common.Hash      `json:"pairHash"`
	QuantityFilled   BigNum           `json:"quantityFilled"`
	Direction        Direction        `json:"direction"`
	Status           OrderStatus      `json:"status"`
}

// EvmSyncStatus contains sync status of EVM state,
// to avoid re-submitting update events from older blocks.
type EvmSyncStatus struct {
	LatestBlockSynced uint64 `json:"latestBlock"`
}

// implement fmt.Stringer
func (m Order) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Order: %v
Status: %s`, m.Order, m.Status))
}

// implement fmt.Stringer
func (m DerivativeOrder) String() string {
	return strings.TrimSpace(fmt.Sprintf(`DerivativeOrder: %v
Status: %s`, m.Order, m.Status))
}

type OrderStatus uint8

const (
	StatusUnknown       OrderStatus = 0
	StatusUnfilled      OrderStatus = 1
	StatusSoftCancelled OrderStatus = 2
	StatusPartialFilled OrderStatus = 3
	StatusFilled        OrderStatus = 4
	StatusExpired       OrderStatus = 5
	StatusHardCancelled OrderStatus = 6
	StatusUntriggered   OrderStatus = 7
)

type Direction uint8

const (
	Long  Direction = 0
	Short Direction = 1
)

func DirectionFromString(direction string) Direction {
	switch direction {
	case "long":
		return Long
	case "short":
		return Short
	default:
		return Long
	}

}

func (d Direction) String() string {
	switch d {
	case Long:
		return "long"
	case Short:
		return "short"
	}
	return ""
}

func OrderStatusFromString(status string) OrderStatus {
	switch status {
	case "unfilled":
		return StatusUnfilled
	case "softCancelled":
		return StatusSoftCancelled
	case "hardCancelled":
		return StatusHardCancelled
	case "partialFilled":
		return StatusPartialFilled
	case "filled":
		return StatusFilled
	case "expired":
		return StatusExpired
	case "untriggered":
		return StatusUntriggered
	default:
		return StatusUnknown
	}
}

func (s OrderStatus) String() string {
	switch s {
	case StatusUnfilled:
		return "unfilled"
	case StatusSoftCancelled:
		return "softCancelled"
	case StatusHardCancelled:
		return "hardCancelled"
	case StatusPartialFilled:
		return "partialFilled"
	case StatusFilled:
		return "filled"
	case StatusExpired:
		return "expired"
	case StatusUntriggered:
		return "untriggered"
	case StatusUnknown:
		return "unknown"
	}
	return ""
}

type OrderFillRequest struct {
	TxHash               common.Hash `json:"txHash"`
	OrderHash            common.Hash `json:"orderHash"`
	ApprovalSignatures   [][]byte    `json:"approvalSignatures"`
	ExpiryAt             int64       `json:"expiryAt"`
	TakerAssetFillAmount BigNum      `json:"takerAssetFillAmount"`
}

type OrderSoftCancelRequest struct {
	TxHash             common.Hash `json:"txHash"`
	OrderHash          common.Hash `json:"orderHash"`
	ApprovalSignatures [][]byte    `json:"approvalSignatures"`
}

type SignedTransaction struct {
	Salt                  BigNum  `json:"salt"`
	SignerAddress         Address `json:"signerAddress"`
	Data                  []byte  `json:"data"`
	ExpirationTimeSeconds BigNum  `json:"expirationTimeSeconds"`
	GasPrice              BigNum  `json:"gasPrice"`

	Domain CoordinatorDomain `json:"domain"`

	Signature []byte `json:"signature"`
}

type CoordinatorDomain struct {
	VerifyingContract Address `json:"verifyingContract"`
	ChainID           BigNum  `json:"chainID"`
}

type Hash struct {
	common.Hash
}

func (h Hash) MarshalJSON() ([]byte, error) {
	hex := h.Hash.Hex()
	buf := make([]byte, 0, len(hex)+2)
	buf = append(buf, '"')
	buf = append(buf, hex...)
	buf = append(buf, '"')
	return buf, nil
}

type HexBytes []byte

func (h HexBytes) MarshalJSON() ([]byte, error) {
	hex := common.ToHex(h)
	buf := make([]byte, 0, len(hex)+2)
	buf = append(buf, '"')
	buf = append(buf, hex...)
	buf = append(buf, '"')
	return buf, nil
}

func (h *HexBytes) UnmarshalJSON(src []byte) error {
	if len(src) == 2 {
		return nil
	} else if len(src) < 2 {
		return errors.New("failed to parse: " + string(src))
	}

	*h = HexBytes(common.FromHex(string(src[1 : len(src)-1])))
	return nil
}

func (h HexBytes) String() string {
	return common.ToHex([]byte(h))
}

type Address struct {
	common.Address
}

func (a Address) MarshalJSON() ([]byte, error) {
	hex := a.Address.Hex()
	buf := make([]byte, 0, len(hex)+2)
	buf = append(buf, '"')
	buf = append(buf, hex...)
	buf = append(buf, '"')
	return buf, nil
}

const nullAddressHex = "0x0000000000000000000000000000000000000000"

func (a Address) IsEmpty() bool {
	if a.Hex() == nullAddressHex {
		return true
	}

	return false
}

type BigNum string

func (n BigNum) Int() *big.Int {
	i := new(big.Int)
	i.SetString(string(n), 10)
	return i
}

func NewBigNum(i *big.Int) BigNum {
	if i == nil {
		return "0"
	}
	return BigNum(i.String())
}

func (n BigNum) Decimal() decimal.Decimal {
	if len(n) == 0 || n[0] == '0' {
		return decimal.New(0, 1)
	}

	d, _ := decimal.NewFromString(string(n))
	return d
}

// TradePair specifies a market of assets exchange.
type TradePair struct {
	Name           string   `json:"name"`
	MakerAssetData HexBytes `json:"makerAssetData"`
	TakerAssetData HexBytes `json:"takerAssetData"`
	TradePairHash  Hash     `json:"hash"`
	Enabled        bool     `json:"enabled"`
}

func (p TradePair) Hash() (common.Hash, error) {
	if (p.TradePairHash.Hash != common.Hash{}) {
		return p.TradePairHash.Hash, nil
	}

	if len(p.MakerAssetData) == 0 {
		return common.Hash{}, errors.New("hash error: no maker asset data specified")
	} else if len(p.TakerAssetData) == 0 {
		return common.Hash{}, errors.New("hash error: no taker asset data specified")
	}

	var hash common.Hash
	if bytes.Compare(p.MakerAssetData, p.TakerAssetData) < 0 {
		hash = common.BytesToHash(keccak256(p.MakerAssetData, p.TakerAssetData))
	} else {
		hash = common.BytesToHash(keccak256(p.TakerAssetData, p.MakerAssetData))
	}

	return hash, nil
}

// implement fmt.Stringer
func (p TradePair) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Pair: %s Enabled: %t`, p.Name, p.Enabled))
}

// DerivativeMarket specifies a derivative market.
type DerivativeMarket struct {
	Ticker       string   `json:"ticker"`
	Oracle       HexBytes `json:"oracle"`
	BaseCurrency HexBytes `json:"baseCurrency"`
	Nonce        BigNum   `json:"nonce"`
	MarketID     Hash     `json:"marketID"`
	Enabled      bool     `json:"enabled"`
}

func (p DerivativeMarket) Hash() (common.Hash, error) {
	if (p.MarketID.Hash != common.Hash{}) {
		return p.MarketID.Hash, nil
	}
	if len(p.Ticker) == 0 {
		return common.Hash{}, errors.New("hash error: no ticker specified")
	} else if len(p.Oracle) == 0 {
		return common.Hash{}, errors.New("hash error: no oracle specified")
	} else if len(p.BaseCurrency) == 0 {
		return common.Hash{}, errors.New("hash error: no BaseCurrency specified")
	} else if len(p.Nonce) == 0 {
		return common.Hash{}, errors.New("hash error: no nonce specified")
	}
	var hash common.Hash

	hash = crypto.Keccak256Hash([]byte(p.Ticker), p.Oracle, p.BaseCurrency, common.BigToHash(p.Nonce.Int()).Bytes())

	//var buf []byte
	//hash := sha3.NewKeccak256()
	//hash.Write(bytes)
	//buf = hash.Sum(buf)
	//
	//log.Println(hexutil.Encode(buf))

	//hash = common.BytesToHash(keccak256([]byte(p.Ticker), p.Oracle, p.BaseCurrency, common.BigToHash(p.Nonce.Int()).Bytes()))
	return hash, nil
}

// implement fmt.Stringer
func (p DerivativeMarket) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Ticker: %s, Oracle %s, BaseCurrency %s, Nonce %s, MarketID: %s, Enabled: %t`, p.Ticker, p.Oracle.String(), p.BaseCurrency.String(), p.Nonce.Decimal().String(), p.MarketID.String(), p.Enabled))
}

// keccak256 calculates and returns the Keccak256 hash of the input data.
func keccak256(data ...[]byte) []byte {
	d := sha3.NewLegacyKeccak256()
	for _, b := range data {
		_, _ = d.Write(b)
	}
	return d.Sum(nil)
}

type ZeroExTransaction struct {
	Type   ZeroExTransactionType `json:"type"`
	Orders []common.Hash         `json:"orders"`
}

type ZeroExTransactionType int

const (
	ZeroExTransactionTypeUnknown   ZeroExTransactionType = 0
	ZeroExOrderFillRequestTx       ZeroExTransactionType = 1
	ZeroExOrderSoftCancelRequestTx ZeroExTransactionType = 2
)
