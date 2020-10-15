package provider

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid"
)

// contextWithCloseChan returns a cancellable context that cancels
// when cancelC chan is closed or when a closeC chan signal arrives.
func contextWithCloseChan(ctx context.Context, closeC <-chan struct{}) (context.Context, func()) {
	cancelC := make(chan struct{})
	closeFn := func() {
		close(cancelC)
	}
	ctx, cancelFn := context.WithCancel(ctx)
	go func(cancelFn func()) {
		select {
		case <-closeC:
			cancelFn()
		case <-cancelC:
			cancelFn()
		}
	}(cancelFn)
	return ctx, closeFn
}

func NewSessionID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).String()
}

func NewSessionIDBytes() []byte {
	id, _ := ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).MarshalText()
	return id
}

func WriteSessionID(dst []byte) ([]byte, error) {
	err := ulid.MustNew(ulid.Timestamp(time.Now()), globalRand).MarshalTextTo(dst)
	return dst, err
}

var globalRand = rand.New(&lockedSource{
	src: rand.NewSource(time.Now().UnixNano()),
})

type lockedSource struct {
	lk  sync.Mutex
	src rand.Source
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}
