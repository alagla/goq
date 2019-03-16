package dispatcher

import (
	"sync"
	"time"
)

// lock with timeout

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

// modified wait group: after Shoot (equivalent to Done) it is released and armed (Add(1))
// again atomically
type ShooterWG struct {
	hold    sync.WaitGroup
	release sync.WaitGroup
}

func NewShooterWG() *ShooterWG {
	return &ShooterWG{}
}

func (sh *ShooterWG) Arm() {
	sh.hold.Add(1)
	sh.release.Add(1)
}

func (sh *ShooterWG) Disarm() {
	sh.hold.Done()
	sh.release.Done()
}

func (sh *ShooterWG) Wait() {
	sh.hold.Wait()
	sh.release.Wait()
}

func (sh *ShooterWG) Shoot() {
	sh.hold.Done()
	sh.hold.Add(1)
	sh.release.Done()
	sh.release.Add(1)
}
