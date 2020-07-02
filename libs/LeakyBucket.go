package libs

import (
	"github.com/porschemacan/golden/libs/internal/clock"
	"sync"
	"sync/atomic"
	"time"
)

type LeakyBucket struct {
	sync.Mutex
	last     time.Time
	sleepFor time.Duration
	rpsFunc  LimitAndWaitFunc
	waitCnt  int32
	clock    Clock
}

type Option func(l *LeakyBucket)

func WithClock(clock Clock) Option {
	return func(l *LeakyBucket) {
		l.clock = clock
	}
}

func NewLeakyBucket(rps LimitAndWaitFunc, opts ...Option) Limiter {
	l := &LeakyBucket{
		rpsFunc: rps,
	}
	for _, opt := range opts {
		opt(l)
	}
	if l.clock == nil {
		l.clock = clock.New()
	}

	return l
}

func (l *LeakyBucket) Take() bool {
	rps, maxWait := l.rpsFunc()
	if rps <= 0 {
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

	// If this is our first request, then we allow it.
	cur := l.clock.Now()
	if l.last.IsZero() {
		l.last = cur
		return true
	}

	perRequest := time.Second / time.Duration(rps)

	// sleepFor calculates how much time we should sleep based on
	// the perRequest budget and how long the last request took.
	// Since the request may take longer than the budget, this number
	// can get negative, and is summed across requests.
	l.sleepFor += perRequest - cur.Sub(l.last)
	l.last = cur

	// We shouldn't allow sleepFor to get too negative, since it would mean that
	// a service that slowed down a lot for a short period of time would get
	// a much higher RPS following that.
	if l.sleepFor < -time.Second {
		l.sleepFor = -time.Second
	}

	// If sleepFor is positive, then we should sleep now.
	if l.sleepFor > 0 {
		l.clock.Sleep(l.sleepFor)
		l.last = cur.Add(l.sleepFor)
		l.sleepFor = 0
	}

	return true
}

func (l *LeakyBucket) Put() {
	//
}

type unlimited struct{}

// NewUnlimited returns a RateLimiter that is not limited.
func NewUnlimited() Limiter {
	return unlimited{}
}

func (unlimited) Take() bool {
	return true
}

func (unlimited) Put() {
}
