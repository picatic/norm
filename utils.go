package norm

import (
	"sync/atomic"
	"sync"
)

// Once is an object that will perform exactly one action.
type OnceDone struct {
	m    sync.Mutex
	done uint32
}

func (o *OnceDone) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 1 {
		return
	}
	// Slow-path.
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

// Return
// added Done to the sync.Once implementation
//
func (o *OnceDone) Done() bool {
	return atomic.LoadUint32(&o.done) == 1
}


