package field

import (
	"sync"
	"sync/atomic"
)

// ShadowInit Making a shadowValue only settable once
type ShadowInit struct {
	m    sync.Mutex
	done uint32
}

// DoInit set a value by func
func (o *ShadowInit) DoInit(f func()) {
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

// InitDone return a value
func (o *ShadowInit) InitDone() bool {
	return atomic.LoadUint32(&o.done) == 1
}
