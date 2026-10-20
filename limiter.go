package main

import (
	"sync"
	"time"
)

type visitor struct {
	requests  int
	windowEnd time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	limit    int
	window   time.Duration
	visitors map[string]*visitor
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		window:   window,
		visitors: make(map[string]*visitor),
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	allowed, _ := rl.Check(key)
	return allowed
}

func (rl *RateLimiter) Check(key string) (bool, time.Duration) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	v, exists := rl.visitors[key]

	if !exists || !now.Before(v.windowEnd) {
		rl.visitors[key] = &visitor{
			requests:  1,
			windowEnd: now.Add(rl.window),
		}
		return true, 0
	}

	if v.requests >= rl.limit {
		return false, time.Until(v.windowEnd)
	}

	v.requests++
	return true, 0
}
