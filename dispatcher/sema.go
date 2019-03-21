package dispatcher

import (
	"time"
)

// lock with timeout

type Sema struct {
	ch chan struct{}
}

func NewSema() *Sema {
	ret := &Sema{
		ch: make(chan struct{}, 1),
	}
	return ret
}

func (cl *Sema) Dispose() {
	close(cl.ch)
}

func (cl *Sema) Acquire(timeout time.Duration) bool {
	if timeout < 0 {
		cl.ch <- struct{}{}
		return true
	}
	select {
	case cl.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (cl *Sema) Release() bool {
	select {
	case <-cl.ch:
		return true
	default:
		return false
	}
}
