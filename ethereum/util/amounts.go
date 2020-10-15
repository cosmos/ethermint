package util

import (
	"fmt"
	"math/big"

	"github.com/shopspring/decimal"
)

type Wei decimal.Decimal

func (w Wei) String() string {
	return (decimal.Decimal)(w).String()
}

func (w Wei) StringGwei() string {
	d := (decimal.Decimal)(w).Div(decimal.NewFromFloat(1e9))
	return d.String()
}

func (w Wei) Bytes() []byte {
	return []byte((decimal.Decimal)(w).String())
}

func (w *Wei) Scan(v interface{}) error {
	if v == nil {
		return nil
	}
	var source string
	switch v.(type) {
	case string:
		source = v.(string)
	case []byte:
		source = string(v.([]byte))
	default:
		return fmt.Errorf("incompatible type for decimal.Decimal: %T", v)
	}
	d, err := decimal.NewFromString(source)
	if err != nil {
		err := fmt.Errorf("failed to parse decimal.Decimal from %s, error: %v", source, err)
		return err
	} else {
		*w = (Wei)(d)
	}
	return nil
}

func (w *Wei) Ether() float64 {
	if w == nil {
		return 0
	}
	f, _ := (*decimal.Decimal)(w).Div(decimal.NewFromFloat(1e18)).Float64()
	return f
}

func (w *Wei) Tokens() float64 {
	return w.Ether()
}

func (w *Wei) ToInt() *big.Int {
	i := big.NewInt(0)
	i.SetString((*decimal.Decimal)(w).String(), 10)
	return i
}

// Gwei is an unsafe way to represent Wei as uint64, used for
// gas price reporting and should not be used for math.
func (w *Wei) Gwei() uint64 {
	if w == nil {
		return 0
	}
	r := big.NewInt(0).Set(w.ToInt())
	m := big.NewInt(0).SetUint64(1e9)
	return r.Div(r, m).Uint64()
}

func Gwei(gwei uint64) *Wei {
	w := big.NewInt(0).SetUint64(gwei)
	m := big.NewInt(0).SetUint64(1e9)
	return BigWei(w.Mul(w, m))
}

func (w *Wei) Mul(m int64) *Wei {
	d := (*decimal.Decimal)(w)
	result := d.Mul(decimal.NewFromBigInt(big.NewInt(m), 0))
	return (*Wei)(&result)
}

func (w *Wei) Div(m int64) *Wei {
	d := (*decimal.Decimal)(w)
	result := d.Div(decimal.NewFromBigInt(big.NewInt(m), 0))
	return (*Wei)(&result)
}

// ToWei converts ether or tokens amount into Wei amount.
func ToWei(amount float64) *Wei {
	d := decimal.NewFromFloat(amount).Mul(decimal.NewFromFloat(1e18))
	return (*Wei)(&d)
}

func DecimalToWei(d decimal.Decimal) *Wei {
	d = d.Mul(decimal.NewFromFloat(1e18))
	return (*Wei)(&d)
}

func DecimalWei(d decimal.Decimal) *Wei {
	return (*Wei)(&d)
}

func BigWei(w *big.Int) *Wei {
	d := decimal.NewFromBigInt(w, 0)
	return (*Wei)(&d)
}

func StringWei(str string) *Wei {
	d, err := decimal.NewFromString(str)
	if err != nil {
		return nil
	}
	return (*Wei)(&d)
}

// Add adds two amounts together and returns a new amount.
func (w *Wei) Add(amount *Wei) *Wei {
	d1 := (*decimal.Decimal)(w)
	d2 := (*decimal.Decimal)(amount)
	result := d1.Add(*d2)
	return (*Wei)(&result)
}

// Sub substracts two amounts and returns a new amount.
func (w *Wei) Sub(amount *Wei) *Wei {
	d1 := (*decimal.Decimal)(w)
	d2 := (*decimal.Decimal)(amount)
	result := d1.Sub(*d2)
	return (*Wei)(&result)
}

// SplitEqual splits the amount into n-1 equal amounts and one remainder.
// Example: (1000).SplitEqual(7) yields [142 142 142 142 142 142 148]
func (w *Wei) SplitEqual(parts int) []*Wei {
	if parts == 0 {
		return nil
	}
	total := big.NewInt(0).Set(w.ToInt())
	q, m := big.NewInt(0).DivMod(total, big.NewInt(int64(parts)), big.NewInt(0))
	if q.Cmp(big.NewInt(0)) == 0 {
		result := make([]*Wei, parts)
		for i := 0; i < parts; i++ {
			result[i] = ToWei(0)
		}
		result[parts-1] = BigWei(m)
		return result
	}
	result := make([]*Wei, 0, parts)
	for total.Cmp(q) >= 0 {
		total.Sub(total, q)
		amount := BigWei(big.NewInt(0).Set(q))
		result = append(result, amount)
	}
	result[len(result)-1] = result[len(result)-1].Add(BigWei(m))
	return result
}
