package dispatcher

import "time"

type AsyncLock struct {
	chIn  chan struct{}
	chOut chan struct{}
}

func NewAsyncLock() *AsyncLock {
	ret := &AsyncLock{
		chIn:  make(chan struct{}),
		chOut: make(chan struct{}),
	}
	go ret.loop()
	ret.chIn <- struct{}{} // initially lock is released
	return ret
}

func (cl *AsyncLock) loop() {
	for range cl.chIn {
		cl.chOut <- struct{}{}
	}
	close(cl.chOut)
}

func (cl *AsyncLock) Destroy() {
	close(cl.chIn)
}

func (cl *AsyncLock) Request(timeout time.Duration) bool {
	if timeout < 0 {
		<-cl.chOut
		return true
	}
	select {
	case <-cl.chOut:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (cl *AsyncLock) Release() bool {
	select {
	case cl.chIn <- struct{}{}:
		return true
	case <-time.After(1 * time.Millisecond):
		// default timeout is 1 milisecond, after that lock is considered unlocked
		return false
	}
}
