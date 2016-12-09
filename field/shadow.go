package field

// ShadowInit Making a shadowValue only settable once
type ShadowInit struct {
	done bool
}

// DoInit set a value by func
func (o *ShadowInit) DoInit(f func()) {
	if o.done == false {
		f()
		o.done = true
	}
}

// InitDone return a value
func (o *ShadowInit) InitDone() bool {
	return o.done
}
