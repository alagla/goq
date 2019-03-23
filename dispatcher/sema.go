package dispatcher

import (
	"time"
)

// lock with timeout

type sema struct {
	ch chan struct{}
}

func newSema() *sema {
	ret := &sema{
		ch: make(chan struct{}, 1),
	}
	return ret
}

func (cl *sema) dispose() {
	close(cl.ch)
}

func (cl *sema) acquire(timeout time.Duration) bool {
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

func (cl *sema) release() bool {
	select {
	case <-cl.ch:
		return true
	default:
		return false
	}
}
