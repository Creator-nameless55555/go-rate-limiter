package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAllowsRequestsWithinLimit(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		if !limiter.Allow("client-1") {
			t.Fatalf("request %d should have been allowed", i+1)
		}
	}
}

func TestRejectsRequestsOverLimit(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	if !limiter.Allow("client-1") {
		t.Fatal("first request should be allowed")
	}
	if !limiter.Allow("client-1") {
		t.Fatal("second request should be allowed")
	}
	if limiter.Allow("client-1") {
		t.Fatal("third request should be rejected")
	}
}

func TestSeparatesClients(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute)

	if !limiter.Allow("client-1") {
		t.Fatal("first client should be allowed")
	}
	if limiter.Allow("client-1") {
		t.Fatal("second request from first client should be rejected")
	}
	if !limiter.Allow("client-2") {
		t.Fatal("different client should be allowed")
	}
}

func TestResetsWindow(t *testing.T) {
	limiter := NewRateLimiter(1, 30*time.Millisecond)

	if !limiter.Allow("client-1") {
		t.Fatal("first request should be allowed")
	}
	if limiter.Allow("client-1") {
		t.Fatal("second request should be rejected within the window")
	}

	time.Sleep(50 * time.Millisecond)

	if !limiter.Allow("client-1") {
		t.Fatal("request should be allowed after the window expires")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello!"))
	})

	handler := rateLimitMiddleware(limiter, mux)

	for i := 1; i <= 3; i++ {
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		req.RemoteAddr = "192.0.2.1:12345"

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)

		if i <= 2 && rec.Code != http.StatusOK {
			t.Errorf("request %d: expected status %d, got %d",
				i, http.StatusOK, rec.Code)
		}

		if i == 3 && rec.Code != http.StatusTooManyRequests {
			t.Errorf("request %d: expected status %d, got %d",
				i, http.StatusTooManyRequests, rec.Code)
		}
	}
}

func TestRateLimitMiddlewareRetryAfterAndJSON(t *testing.T) {
	limiter := NewRateLimiter(1, time.Minute)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := rateLimitMiddleware(limiter, next)

	makeRequest := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		req.RemoteAddr = "192.0.2.10:12345"

		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, req)
		return rec
	}

	first := makeRequest()
	if first.Code != http.StatusOK {
		t.Fatalf("first request: expected status %d, got %d",
			http.StatusOK, first.Code)
	}

	second := makeRequest()
	if second.Code != http.StatusTooManyRequests {
		t.Fatalf("second request: expected status %d, got %d",
			http.StatusTooManyRequests, second.Code)
	}

	retryAfter := second.Header().Get("Retry-After")
	if retryAfter == "" {
		t.Fatal("expected Retry-After header to be set")
	}

	if !strings.Contains(second.Header().Get("Content-Type"), "application/json") {
		t.Fatalf("expected JSON Content-Type, got %q",
			second.Header().Get("Content-Type"))
	}

	var body struct {
		Error      string `json:"error"`
		RetryAfter int    `json:"retry_after"`
	}

	if err := json.Unmarshal(second.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if body.Error != "rate limit exceeded" {
		t.Errorf("unexpected error message: %q", body.Error)
	}

	if body.RetryAfter < 1 {
		t.Errorf("expected retry_after >= 1, got %d", body.RetryAfter)
	}
}
