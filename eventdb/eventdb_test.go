package eventdb

import (
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestAddEvent(t *testing.T) {
	require := require.New(t)
	noHash := common.Hash{}

	db := NewOrderEventDB()

	ev := &OrderEvent{
		Type:       OrderUpdateFilled,
		BlockNum:   1,
		TxHash:     noHash,
		OrderHash:  noHash,
		FillAmount: big.NewInt(0),
	}
	db.AddEvent(ev)

	res, ok := db.GetFillEvent(1, noHash, noHash)

	require.True(ok)
	require.Equal(ev, res)
}

func TestGetFillEvent(t *testing.T) {
	require := require.New(t)
	noHash := common.Hash{}

	db := NewOrderEventDB()

	ev := &OrderEvent{
		Type:       OrderUpdateFilled,
		BlockNum:   1,
		TxHash:     noHash,
		OrderHash:  noHash,
		FillAmount: big.NewInt(0),
	}
	db.AddEvent(ev)

	res, ok := db.GetFillEvent(1, noHash, noHash)

	require.True(ok)
	require.Equal(ev, res)

	_, ok = db.GetCancelEvent(1, noHash, noHash)
	require.False(ok)
}

func TestGetCancelEvent(t *testing.T) {
	require := require.New(t)
	noHash := common.Hash{}

	db := NewOrderEventDB()

	ev := &OrderEvent{
		Type:       OrderUpdateHardCancelled,
		BlockNum:   1,
		TxHash:     noHash,
		OrderHash:  noHash,
		FillAmount: big.NewInt(0),
	}
	db.AddEvent(ev)

	res, ok := db.GetCancelEvent(1, noHash, noHash)

	require.True(ok)
	require.Equal(ev, res)

	_, ok = db.GetFillEvent(1, noHash, noHash)
	require.False(ok)
}

func TestRangeFillEvents(t *testing.T) {
	require := require.New(t)

	db := NewOrderEventDB()

	events := make([]*OrderEvent, 5)
	for i := 0; i < 5; i++ {
		events[i] = &OrderEvent{
			Type:       OrderUpdateFilled,
			BlockNum:   uint64(i),
			TxHash:     randomHash(),
			OrderHash:  randomHash(),
			FillAmount: big.NewInt(0),
		}
		db.AddEvent(events[i])
	}

	var seen int
	db.RangeFillEvents(4, func(ev *OrderEvent) error {
		require.Equal(events[ev.BlockNum], ev)
		require.LessOrEqual(ev.BlockNum, uint64(3))
		seen++
		return nil
	})
	require.Equal(4, seen)
}

func TestRangeCancelEvents(t *testing.T) {
	require := require.New(t)

	db := NewOrderEventDB()

	events := make([]*OrderEvent, 5)
	for i := 0; i < 5; i++ {
		events[i] = &OrderEvent{
			Type:       OrderUpdateHardCancelled,
			BlockNum:   uint64(i),
			TxHash:     randomHash(),
			OrderHash:  randomHash(),
			FillAmount: big.NewInt(0),
		}
		db.AddEvent(events[i])
	}

	var seen int
	db.RangeCancelEvents(4, func(ev *OrderEvent) error {
		require.Equal(events[ev.BlockNum], ev)
		require.LessOrEqual(ev.BlockNum, uint64(3))
		seen++
		return nil
	})
	require.Equal(4, seen)
}

func TestForgetFillEvents(t *testing.T) {
	require := require.New(t)

	db := NewOrderEventDB()

	events := make([]*OrderEvent, 5)
	for i := 0; i < 5; i++ {
		events[i] = &OrderEvent{
			Type:       OrderUpdateFilled,
			BlockNum:   uint64(i),
			TxHash:     randomHash(),
			OrderHash:  randomHash(),
			FillAmount: big.NewInt(0),
		}
		db.AddEvent(events[i])
	}

	db.ForgetFillEvents(4)

	_, ok := db.GetFillEvent(events[4].BlockNum, events[4].TxHash, events[4].OrderHash)
	require.True(ok)

	for i := 3; i >= 0; i-- {
		_, ok := db.GetFillEvent(events[i].BlockNum, events[i].TxHash, events[i].OrderHash)
		require.False(ok)
	}
}

func TestForgetCancelEvents(t *testing.T) {
	require := require.New(t)

	db := NewOrderEventDB()

	events := make([]*OrderEvent, 5)
	for i := 0; i < 5; i++ {
		events[i] = &OrderEvent{
			Type:       OrderUpdateHardCancelled,
			BlockNum:   uint64(i),
			TxHash:     randomHash(),
			OrderHash:  randomHash(),
			FillAmount: big.NewInt(0),
		}
		db.AddEvent(events[i])
	}

	db.ForgetCancelEvents(4)

	_, ok := db.GetCancelEvent(events[4].BlockNum, events[4].TxHash, events[4].OrderHash)
	require.True(ok)

	for i := 3; i >= 0; i-- {
		_, ok := db.GetCancelEvent(events[i].BlockNum, events[i].TxHash, events[i].OrderHash)
		require.False(ok)
	}
}

func TestSetCurrentBlock(t *testing.T) {
	require := require.New(t)

	db := NewOrderEventDB()

	wg := new(sync.WaitGroup)
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()

			db.SetCurrentBlock(uint64(i))
		}(i)
	}
	wg.Wait()

	db.SetCurrentBlock(999)
	require.Equal(uint64(999), db.CurrentBlock())
}

func TestCurrentBlock(t *testing.T) {
	require := require.New(t)

	db := NewOrderEventDB()
	db.SetCurrentBlock(5)

	require.Equal(uint64(5), db.CurrentBlock())
}

func randomHash() common.Hash {
	buf := make([]byte, 32)
	_, _ = rand.Read(buf)
	return common.BytesToHash(buf)
}
