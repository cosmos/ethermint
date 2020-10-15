package eventdb

import (
	"bytes"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/xlab/suplog"
	btree "modernc.org/b"
)

type FuturesPositionEventType int

const (
	FuturesPositionUpdateFilled        FuturesPositionEventType = 1
	FuturesPositionUpdateHardCancelled FuturesPositionEventType = 2
)

type FuturesPositionEvent struct {
	Type           FuturesPositionEventType `json:"type"`
	BlockNum       uint64                   `json:"blockNum"`
	TxHash         common.Hash              `json:"txHash"`
	MakerAddress   common.Address           `json:"makerAddress"`
	AccountID      common.Hash              `json:"accountID"`
	OrderHash      common.Hash              `json:"orderHash"`
	MarketID       common.Hash              `json:"marketId"`
	ContractPrice  *big.Int                 `json:"contractPrice"`
	QuantityFilled *big.Int                 `json:"quantityFilled"`
	PositionID     *big.Int                 `json:"positionId"`
	IsLong         bool                     `json:"isLong"`

	key []byte
}

func (e *FuturesPositionEvent) Equals(e2 *FuturesPositionEvent) bool {
	if e == nil && e2 == nil {
		return true
	} else if e == nil || e2 == nil {
		return false
	}

	if e.IsLong != e2.IsLong {
		return false
	}

	if e.Type != e2.Type {
		return false
	} else if e.BlockNum != e2.BlockNum {
		return false
	} else if e.TxHash != e2.TxHash {
		return false
	}

	if !isEqualInt(e.QuantityFilled, e2.QuantityFilled) {
		return false
	} else if !isEqualInt(e.ContractPrice, e2.ContractPrice) {
		return false
	} else if !isEqualInt(e.PositionID, e2.PositionID) {
		return false
	}

	if e.MakerAddress != e2.MakerAddress {
		return false
	} else if e.OrderHash != e2.OrderHash {
		return false
	} else if e.MarketID != e2.MarketID {
		return false
	}

	return true
}

func isEqualInt(a, b *big.Int) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	}

	if a.Cmp(b) != 0 {
		return false
	}

	return true
}

func (e *FuturesPositionEvent) Key() []byte {
	if e.key != nil {
		return e.key
	}

	var tag = 0
	if e.IsLong {
		tag = 1
	}

	keyStr := fmt.Sprintf("%s_%s_%s_%d", padBlockNum(e.BlockNum), e.TxHash.Hex(), e.OrderHash.Hex(), tag)
	e.key = []byte(keyStr)
	return e.key
}

func (_ *FuturesPositionEvent) CmpFunc() func(a, b interface{}) int {
	return func(a, b interface{}) int {
		return bytes.Compare(a.([]byte), b.([]byte))
	}
}

type FuturesPositionRangeFunc func(ev *FuturesPositionEvent) error

type FuturesPositionEventDB interface {
	AddEvent(ev *FuturesPositionEvent)
	SetCurrentBlock(blockNum uint64)
	CurrentBlock() uint64
	GetFillEvent(blockNum uint64, txHash, orderHash common.Hash, isLong bool) (*FuturesPositionEvent, bool)
	GetCancelEvent(blockNum uint64, txHash, orderHash common.Hash) (*FuturesPositionEvent, bool)
	RangeFillEvents(from uint64, fn FuturesPositionRangeFunc) error
	RangeCancelEvents(from uint64, fn FuturesPositionRangeFunc) error
	ForgetFillEvents(from uint64)
	ForgetCancelEvents(from uint64)
}

func NewFuturesPositionEventDB() FuturesPositionEventDB {
	return &memFuturesPositionEventDB{
		blockNum:     0,
		fillEvents:   btree.TreeNew((&FuturesPositionEvent{}).CmpFunc()),
		cancelEvents: btree.TreeNew((&FuturesPositionEvent{}).CmpFunc()),
		blockMux:     new(sync.RWMutex),
		fillMux:      new(sync.RWMutex),
		cancelMux:    new(sync.RWMutex),
	}
}

type memFuturesPositionEventDB struct {
	blockNum uint64
	blockMux *sync.RWMutex

	fillEvents *btree.Tree
	fillMux    *sync.RWMutex

	cancelEvents *btree.Tree
	cancelMux    *sync.RWMutex
}

func (db *memFuturesPositionEventDB) AddEvent(ev *FuturesPositionEvent) {
	switch ev.Type {
	case FuturesPositionUpdateFilled:
		db.fillMux.Lock()
		defer db.fillMux.Unlock()

		db.fillEvents.Set(ev.Key(), ev)

	case FuturesPositionUpdateHardCancelled:
		db.cancelMux.Lock()
		defer db.cancelMux.Unlock()

		db.cancelEvents.Set(ev.Key(), ev)
	default:
		log.WithField("type", ev.Type).Errorln("unsupported event not inserted in DB")
		return
	}
}

func (db *memFuturesPositionEventDB) GetFillEvent(blockNum uint64, txHash, orderHash common.Hash, isLong bool) (*FuturesPositionEvent, bool) {
	db.fillMux.RLock()
	defer db.fillMux.RUnlock()

	var tag = 0
	if isLong {
		tag = 1
	}

	keyStr := fmt.Sprintf("%s_%s_%s_%d", padBlockNum(blockNum), txHash.Hex(), orderHash.Hex(), tag)

	ev, ok := db.fillEvents.Get([]byte(keyStr))
	if !ok {
		return nil, false
	}

	return ev.(*FuturesPositionEvent), true
}

func (db *memFuturesPositionEventDB) GetCancelEvent(blockNum uint64, txHash, orderHash common.Hash) (*FuturesPositionEvent, bool) {
	db.cancelMux.RLock()
	defer db.cancelMux.RUnlock()

	keyStr := fmt.Sprintf("%s_%s_%s_0", padBlockNum(blockNum), txHash.Hex(), orderHash.Hex())

	ev, ok := db.cancelEvents.Get([]byte(keyStr))
	if !ok {
		return nil, false
	}

	return ev.(*FuturesPositionEvent), true
}

func (db *memFuturesPositionEventDB) RangeFillEvents(from uint64, fn FuturesPositionRangeFunc) error {
	db.fillMux.RLock()
	defer db.fillMux.RUnlock()

	prefix := []byte(fmt.Sprintf("%s_", padBlockNum(from)))

	cur, _ := db.fillEvents.Seek(prefix)
	defer cur.Close()

	for {
		k, v, err := cur.Prev()
		if err != nil {
			return nil
		}

		if err := fn(v.(*FuturesPositionEvent)); err != nil {
			if err == ErrRangeStop {
				return nil
			}

			log.WithError(err).Warningf("error processing item %s", string(k.([]byte)))
			return err
		}
	}
}

func (db *memFuturesPositionEventDB) RangeCancelEvents(from uint64, fn FuturesPositionRangeFunc) error {
	db.cancelMux.RLock()
	defer db.cancelMux.RUnlock()

	prefix := []byte(fmt.Sprintf("%s_", padBlockNum(from)))

	cur, _ := db.cancelEvents.Seek(prefix)
	defer cur.Close()

	for {
		k, v, err := cur.Prev()
		if err != nil {
			return nil
		}

		if err := fn(v.(*FuturesPositionEvent)); err != nil {
			if err == ErrRangeStop {
				return nil
			}

			log.WithError(err).Warningf("error processing item %s", string(k.([]byte)))
			return err
		}
	}
}

func (db *memFuturesPositionEventDB) ForgetFillEvents(from uint64) {
	db.fillMux.Lock()
	defer db.fillMux.Unlock()

	prefix := []byte(fmt.Sprintf("%s_", padBlockNum(from)))

	cur, _ := db.fillEvents.Seek(prefix)
	defer cur.Close()

	for {
		k, _, err := cur.Prev()
		if err != nil {
			return
		}

		_ = db.fillEvents.Delete(k)
	}
}

func (db *memFuturesPositionEventDB) ForgetCancelEvents(from uint64) {
	db.cancelMux.Lock()
	defer db.cancelMux.Unlock()

	prefix := []byte(fmt.Sprintf("%s_", padBlockNum(from)))

	cur, _ := db.cancelEvents.Seek(prefix)
	defer cur.Close()

	for {
		k, _, err := cur.Prev()
		if err != nil {
			return
		}

		_ = db.cancelEvents.Delete(k)
	}
}

func (db *memFuturesPositionEventDB) SetCurrentBlock(blockNum uint64) {
	db.blockMux.Lock()
	defer db.blockMux.Unlock()

	db.blockNum = blockNum
}

func (db *memFuturesPositionEventDB) CurrentBlock() uint64 {
	db.blockMux.RLock()
	defer db.blockMux.RUnlock()

	return db.blockNum
}
