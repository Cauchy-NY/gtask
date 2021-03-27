package gtask

import (
	"sync"
	"sync/atomic"
)

// Once is an object that will perform exactly one action
// until Reset is called.
// See http://golang.org/pkg/sync/#Once
type Once struct {
	m    sync.Mutex
	done uint32
}

// Do calls the function f if and only if Do has not been invoked
// without error for this instance of Once.  In other words, given
// 	var once Once
// if once.Do(f) is called multiple times, only the first call will
// invoke f, even if f has a different value in each invocation unless
// f returns an error.  A new instance of Once is required for each
// function to execute.
//
// Do is intended for initialization that must be run exactly once.  Since f
// is niladic, it may be necessary to use a function literal to capture the
// arguments to a function to be invoked by Do:
// 	err := config.once.Do(func() error { return config.init(filename) })
func (o *Once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	if o.done == 0 {
		err = f()
		if err == nil {
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}

// Reset indicates that the next call to Do should actually be called
// once again.
func (o *Once) Reset() {
	o.m.Lock()
	defer o.m.Unlock()
	atomic.StoreUint32(&o.done, 0)
}
