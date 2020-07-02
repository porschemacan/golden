package libs

import (
	"fmt"
	"github.com/porschemacan/golden/libs/internal/clock"
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func ExampleRatelimit() {
	rl := NewLeakyBucket(func() (int32, int32) {
		return 100, 100
	})

	for i := 0; i < 10; i++ {
		ret := rl.Take()
		fmt.Println(i, ret)
	}

}

func TestUnlimited(t *testing.T) {
	now := time.Now()
	rl := NewUnlimited()
	for i := 0; i < 1000; i++ {
		rl.Take()
	}
	assert.Condition(t, func() bool { return time.Now().Sub(now) < 1*time.Millisecond }, "no artificial delay")
}

func TestRateLimiter(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	clock := clock.NewMock()
	rl := NewLeakyBucket(func() (int32, int32) {
		return 100, 100
	}, WithClock(clock))

	count := int32(0)

	// Until we're done...
	done := make(chan struct{})
	defer close(done)

	// Create copious counts concurrently.
	go job(rl, &count, done)
	go job(rl, &count, done)
	go job(rl, &count, done)
	go job(rl, &count, done)

	clock.AfterFunc(1*time.Second, func() {
		assert.InDelta(t, 100, count, 10, "count within rate limit")
	})

	clock.AfterFunc(2*time.Second, func() {
		assert.InDelta(t, 200, count, 10, "count within rate limit")
	})

	clock.AfterFunc(3*time.Second, func() {
		assert.InDelta(t, 300, count, 10, "count within rate limit")
		wg.Done()
	})

	clock.Add(4 * time.Second)

	clock.Add(5 * time.Second)
}

func TestDelayedRateLimiter(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	clock := clock.NewMock()
	slow := NewLeakyBucket(func() (int32, int32) {
		return 10, 10
	}, WithClock(clock))

	fast := NewLeakyBucket(func() (int32, int32) {
		return 100, 100
	}, WithClock(clock))

	count := int32(0)

	// Until we're done...
	done := make(chan struct{})
	defer close(done)

	// Run a slow job
	go func() {
		for {
			slow.Take()
			fast.Take()
			atomic.AddInt32(&count, 1)
			select {
			case <-done:
				return
			default:
			}
		}
	}()

	// Accumulate slack for 10 seconds,
	clock.AfterFunc(20*time.Second, func() {
		// Then start working.
		go job(fast, &count, done)
		go job(fast, &count, done)
		go job(fast, &count, done)
		go job(fast, &count, done)
	})

	clock.AfterFunc(30*time.Second, func() {
		assert.InDelta(t, 1300, count, 10, "count within rate limit")
		wg.Done()
	})

	clock.Add(40 * time.Second)
}

func job(rl Limiter, count *int32, done <-chan struct{}) {
	for {
		rl.Take()
		atomic.AddInt32(count, 1)
		select {
		case <-done:
			return
		default:
		}
	}
}
