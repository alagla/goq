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

func (sem *sema) dispose() {
	close(sem.ch)
}

func (sem *sema) acquire(timeout time.Duration) bool {
	if timeout < 0 {
		sem.ch <- struct{}{}
		return true
	}
	select {
	case sem.ch <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (sem *sema) release() bool {
	select {
	case <-sem.ch:
		return true
	default:
		return false
	}
}
