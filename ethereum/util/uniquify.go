package util

import (
	"fmt"
	"sync"
)

// Uniquify is a type of advanced mutex. It allows to create named resource locks.
type Uniquify interface {
	// Call executes only one callable with same id at a time.
	// Multilpe asynchronous calls with same id will be executed sequentally.
	Call(id string, callable func() error) error
}

// NewUniquify returns a new thread-safe uniquify object.
func NewUniquify() Uniquify {
	return &uniquify{
		tasks: make(map[string]*sync.WaitGroup),
	}
}

type uniquify struct {
	lock  sync.Mutex
	tasks map[string]*sync.WaitGroup
}

func (u *uniquify) Call(id string, callable func() error) error {
	var errC = make(chan error)
	func() {
		u.lock.Lock()
		defer u.lock.Unlock()
		oldWg := u.tasks[id]
		if oldWg != nil {
			go func() {
				oldWg.Wait()
				go func() {
					errC <- u.Call(id, callable)
				}()
			}()
			return
		}
		wg := new(sync.WaitGroup)
		wg.Add(1)
		u.tasks[id] = wg

		go func() {
			var err error
			defer func() {
				errC <- err

				u.lock.Lock()
				defer u.lock.Unlock()
				delete(u.tasks, id)
			}()
			defer wg.Done()
			defer func(err *error) {
				if panicData := recover(); panicData != nil {
					if e, ok := panicData.(error); ok {
						*err = e
						return
					}
					*err = fmt.Errorf("%+v", panicData)
				}
			}(&err)
			err = callable()
		}()
	}()

	return <-errC
}
