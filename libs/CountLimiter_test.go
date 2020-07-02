package libs

import (
	"github.com/porschemacan/golden/libs/internal/clock"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func jobConn(t *clock.Mock, rl Limiter, count *int32, done <-chan struct{}) {
	for {
		if rl.Take() {
			atomic.AddInt32(count, 1)
			t.AfterFunc(1100*time.Millisecond, func() {
				rl.Put()
			})
		}

		select {
		case <-done:
			return
		default:
		}
	}
}

func TestCounterLimiter(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	clock := clock.NewMock()
	rl := NewCountLimiter(func() (int32, int32) {
		return 100, 100
	})

	count := int32(0)

	// Until we're done...
	done := make(chan struct{})
	defer close(done)

	// Create copious counts concurrently.
	go jobConn(clock, rl, &count, done)
	go jobConn(clock, rl, &count, done)
	go jobConn(clock, rl, &count, done)
	go jobConn(clock, rl, &count, done)

	triggerCnt := 0

	clock.AfterFunc(1*time.Second, func() {
		triggerCnt++
		assert.InDelta(t, 100, count, 10, "count within rate limit")
	})

	clock.AfterFunc(2*time.Second, func() {
		triggerCnt++
		assert.InDelta(t, 200, count, 10, "count within rate limit")
	})

	clock.AfterFunc(3*time.Second, func() {
		triggerCnt++
		assert.InDelta(t, 300, count, 10, "count within rate limit")
		wg.Done()
	})

	time.Sleep(100 * time.Millisecond)
	clock.Add(1200 * time.Millisecond)

	time.Sleep(100 * time.Millisecond)
	clock.Add(1200 * time.Millisecond)

	time.Sleep(100 * time.Millisecond)
	clock.Add(1200 * time.Millisecond)

	assert.Equal(t, 3, triggerCnt)
}
