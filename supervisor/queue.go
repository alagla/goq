package supervisor

import (
	"sync"
	"time"
)

// fast fifo queue for supervisor input

const maxQueueLen = 100

type effectQueue struct {
	sync.Mutex
	data  [maxQueueLen]*quantMsg
	first int
	last  int
}

func newEffectQueue() *effectQueue {
	return &effectQueue{}
}

func (que *effectQueue) _put(elem *quantMsg) bool {
	next := (que.last + 1) % maxQueueLen
	if next == que.first {
		return false
	}
	que.data[next] = elem
	return true
}

// puts element into queue. Timeout - how long to wait for free buffer space
func (que *effectQueue) put(elem *quantMsg, timeout time.Duration) bool {
	start := time.Now()
	for time.Since(start) < timeout {
		que.Lock()
		if que._put(elem) {
			que.Unlock()
			return true
		}
		que.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return false // timeout
}

func (que *effectQueue) _get() *quantMsg {
	if que.first == que.last {
		return nil
	}
	ret := que.data[que.first]
	que.first = (que.first + 1) % maxQueueLen
	return ret
}

func (que *effectQueue) poll(timeout time.Duration) *quantMsg {
	start := time.Now()
	var ret *quantMsg
	for time.Since(start) < timeout {
		que.Lock()
		ret = que._get()
		if ret != nil {
			que.Unlock()
			return ret
		}
		que.Unlock()
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}
