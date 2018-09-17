package utils

import "sync"

type SafeFlag struct {
	*sync.RWMutex
	value bool
}

func NewSafeFalg() *SafeFlag {
	return &SafeFlag{
		&sync.RWMutex{}, false,
	}
}


func (x *SafeFlag) SetValue(value bool) {
	f := false
	x.RLock()
	if x.value == value {
		f = true
	}
	x.RUnlock()
	if f {
		return
	}
	x.Lock()
	x.value = value
	x.Unlock()
}

func (x *SafeFlag) Value() (r bool) {
	x.RLock()
	r = x.value
	x.RUnlock()
	return
}
