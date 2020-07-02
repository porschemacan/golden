package libs

import (
	"sync"
	"sync/atomic"
)

func (l *CountLimiter) Take() bool {
	maxCount, maxWait := l.countFunc()
	if maxCount <= 0 {
		return true
	}
	if l.waitCnt >= maxWait {
		return false
	}

	atomic.AddInt32(&l.waitCnt, 1)
	l.Lock()
	defer (func() {
		l.Unlock()
		atomic.AddInt32(&l.waitCnt, -1)
	})()

	if l.curCnt >= maxCount {
		return false
	}

	l.curCnt++
	return true
}

func (l *CountLimiter) Put() {
	l.Lock()
	defer (func() {
		l.Unlock()
	})()

	l.curCnt--
}

func NewCountLimiter(cf LimitAndWaitFunc) Limiter {
	l := &CountLimiter{
		countFunc: cf,
	}

	return l
}

type CountLimiter struct {
	sync.Mutex
	waitCnt   int32
	curCnt    int32
	countFunc LimitAndWaitFunc
}
