package dispatcher

import "time"

type LockWithTimeout struct {
	ch chan struct{}
}

func NewAsyncLock() *LockWithTimeout {
	ret := &LockWithTimeout{
		ch: make(chan struct{}, 1),
	}
	return ret
}

func (cl *LockWithTimeout) Destroy() {
	close(cl.ch)
}

func (cl *LockWithTimeout) Acquire(timeout time.Duration) bool {
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

func (cl *LockWithTimeout) Release() bool {
	select {
	case <-cl.ch:
		return true
	default:
		return false
	}
}
