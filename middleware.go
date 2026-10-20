package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/http"
	"time"
)

type RateLimitError struct {
	Error      string `json:"error"`
	RetryAfter int    `json:"retry_after"`
}

func rateLimitMiddleware(
	limiter *RateLimiter,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "invalid client address", http.StatusBadRequest)
			return
		}

		allowed, wait := limiter.Check(ip)

		if !allowed {
			retryAfter := int(math.Ceil(wait.Seconds()))
			if retryAfter < 1 {
				retryAfter = 1
			}

			w.Header().Set(
				"Retry-After",
				time.Duration(retryAfter).String(),
			)

			// Retry-After expects seconds, not a duration string.
			w.Header().Set("Retry-After", fmt.Sprint(retryAfter))
			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(RateLimitError{
				Error:      "rate limit exceeded",
				RetryAfter: retryAfter,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
