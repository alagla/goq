package dispatcher

import "time"

type AsyncLock struct {
	ch chan struct{}
}

func NewAsyncLock() *AsyncLock {
	ret := &AsyncLock{
		ch: make(chan struct{}, 1),
	}
	return ret
}

func (cl *AsyncLock) Destroy() {
	close(cl.ch)
}

func (cl *AsyncLock) Acquire(timeout time.Duration) bool {
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

func (cl *AsyncLock) Release() bool {
	select {
	case <-cl.ch:
		return true
	default:
		return false
	}
}
