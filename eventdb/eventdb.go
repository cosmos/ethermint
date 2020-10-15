package eventdb

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/xlab/suplog"
	btree "modernc.org/b"
)

type OrderEventType int

const (
	OrderUpdateFilled        OrderEventType = 1
	OrderUpdateHardCancelled OrderEventType = 2
)

type OrderEvent struct {
	Type       OrderEventType `json:"type"`
	BlockNum   uint64         `json:"blockNum"`
	TxHash     common.Hash    `json:"txHash"`
	OrderHash  common.Hash    `json:"orderHash"`
	FillAmount *big.Int       `json:"fillAmount"`

	key []byte
}

func (e *OrderEvent) Equals(e2 *OrderEvent) bool {
	if e == nil && e2 == nil {
		return true
	} else if e == nil || e2 == nil {
		return false
	}

	if e.Type != e2.Type {
		return false
	} else if e.BlockNum != e2.BlockNum {
		return false
	} else if e.TxHash != e2.TxHash {
		return false
	} else if e.OrderHash != e2.OrderHash {
		return false
	}

	if e.FillAmount == nil && e2.FillAmount == nil {
		return true
	} else if e.FillAmount == nil || e2.FillAmount == nil {
		return false
	}

	if e.FillAmount.Cmp(e2.FillAmount) != 0 {
		return false
	}

	return true
}

func (e *OrderEvent) Key() []byte {
	if e.key != nil {
		return e.key
	}

	keyStr := fmt.Sprintf("%s_%s_%s", padBlockNum(e.BlockNum), e.TxHash.Hex(), e.OrderHash.Hex())
	e.key = []byte(keyStr)
	return e.key
}

func (_ *OrderEvent) CmpFunc() func(a, b interface{}) int {
	return func(a, b interface{}) int {
		return bytes.Compare(a.([]byte), b.([]byte))
	}
}

var ErrRangeStop = errors.New("stop")

type OrderRangeFunc func(ev *OrderEvent) error

type OrderEventDB interface {
	AddEvent(ev *OrderEvent)
	SetCurrentBlock(blockNum uint64)
	CurrentBlock() uint64
	GetFillEvent(blockNum uint64, txHash, orderHash common.Hash) (*OrderEvent, bool)
	GetCancelEvent(blockNum uint64, txHash, orderHash common.Hash) (*OrderEvent, bool)
	RangeFillEvents(from uint64, fn OrderRangeFunc) error
	RangeCancelEvents(from uint64, fn OrderRangeFunc) error
	ForgetFillEvents(from uint64)
	ForgetCancelEvents(from uint64)
}

func NewOrderEventDB() OrderEventDB {
	return &memOrderEventDB{
		blockNum:     0,
		fillEvents:   btree.TreeNew((&OrderEvent{}).CmpFunc()),
		cancelEvents: btree.TreeNew((&OrderEvent{}).CmpFunc()),
		blockMux:     new(sync.RWMutex),
		fillMux:      new(sync.RWMutex),
		cancelMux:    new(sync.RWMutex),
	}
}

type memOrderEventDB struct {
	blockNum uint64
	blockMux *sync.RWMutex

	fillEvents *btree.Tree
	fillMux    *sync.RWMutex

	cancelEvents *btree.Tree
	cancelMux    *sync.RWMutex
}

func (db *memOrderEventDB) AddEvent(ev *OrderEvent) {
	switch ev.Type {
	case OrderUpdateFilled:
		db.fillMux.Lock()
		defer db.fillMux.Unlock()

		db.fillEvents.Set(ev.Key(), ev)

	case OrderUpdateHardCancelled:
		db.cancelMux.Lock()
		defer db.cancelMux.Unlock()

		db.cancelEvents.Set(ev.Key(), ev)
	default:
		log.WithField("type", ev.Type).Errorln("unsupported event not inserted in DB")
		return
	}
}

func (db *memOrderEventDB) GetFillEvent(blockNum uint64, txHash, orderHash common.Hash) (*OrderEvent, bool) {
	db.fillMux.RLock()
	defer db.fillMux.RUnlock()

	keyStr := fmt.Sprintf("%s_%s_%s", padBlockNum(blockNum), txHash.Hex(), orderHash.Hex())

	ev, ok := db.fillEvents.Get([]byte(keyStr))
	if !ok {
		return nil, false
	}

	return ev.(*OrderEvent), true
}

func (db *memOrderEventDB) GetCancelEvent(blockNum uint64, txHash, orderHash common.Hash) (*OrderEvent, bool) {
	db.cancelMux.RLock()
	defer db.cancelMux.RUnlock()

	keyStr := fmt.Sprintf("%s_%s_%s", padBlockNum(blockNum), txHash.Hex(), orderHash.Hex())

	ev, ok := db.cancelEvents.Get([]byte(keyStr))
	if !ok {
		return nil, false
	}

	return ev.(*OrderEvent), true
}

func (db *memOrderEventDB) RangeFillEvents(from uint64, fn OrderRangeFunc) error {
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

		if err := fn(v.(*OrderEvent)); err != nil {
			if err == ErrRangeStop {
				return nil
			}

			log.WithField("key", string(k.([]byte))).WithError(err).Errorln("error processing item")
			return err
		}
	}
}

func (db *memOrderEventDB) RangeCancelEvents(from uint64, fn OrderRangeFunc) error {
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

		if err := fn(v.(*OrderEvent)); err != nil {
			if err == ErrRangeStop {
				return nil
			}

			log.WithField("key", string(k.([]byte))).WithError(err).Warningf("error processing item")
			return err
		}
	}
}

func (db *memOrderEventDB) ForgetFillEvents(from uint64) {
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

func (db *memOrderEventDB) ForgetCancelEvents(from uint64) {
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

func (db *memOrderEventDB) SetCurrentBlock(blockNum uint64) {
	db.blockMux.Lock()
	defer db.blockMux.Unlock()

	db.blockNum = blockNum
}

func (db *memOrderEventDB) CurrentBlock() uint64 {
	db.blockMux.RLock()
	defer db.blockMux.RUnlock()

	return db.blockNum
}

// padBlockNum creates a byte-comparable 0-padded string
// from block number. Width is set to 20 digits, which should last very long.
func padBlockNum(num uint64) string {
	return fmt.Sprintf("%020d", num)
}
