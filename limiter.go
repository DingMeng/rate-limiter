package rate_limiter

import (
	"sync"
	"sync/atomic"
	"time"
)

type RateLimiter struct {
	rate      uint64
	lastCheck uint64
	max       uint64
	remaining uint64
	step      uint64
	mu        sync.Mutex
}

func New(rate int, per time.Duration) *RateLimiter {
	nanoPer := uint64(per)
	rate64 := uint64(rate)
	rl := &RateLimiter{
		rate:      rate64,
		step: nanoPer,
		max:       rate64 * nanoPer,
		remaining: rate64 * nanoPer,
		lastCheck: NanoNow(),
	}
	return rl
}

func (rl *RateLimiter) Limit() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := NanoNow()
	rate := atomic.LoadUint64(&rl.rate)
	step := atomic.LoadUint64(&rl.step)
	passed := now - atomic.SwapUint64(&rl.lastCheck, now)
	current := atomic.AddUint64(&rl.remaining , passed*rate)
	if max:=atomic.LoadUint64(&rl.max);current>max{
		atomic.SwapUint64(&rl.remaining,max)
		current = max
	}
	if current < step {
		return true
	}
	atomic.CompareAndSwapUint64(&rl.remaining,current,current-step)
	return false
}

func (rl *RateLimiter) ChangeRate(rate int,per time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	atomic.SwapUint64(&rl.rate,uint64(rate))
	atomic.SwapUint64(&rl.step,uint64(per))
}

func NanoNow() uint64 {
	return uint64(time.Now().UnixNano())
}
