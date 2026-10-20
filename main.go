package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	limiter := NewRateLimiter(5, 10*time.Second)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hello! Your request was accepted.\n"))
	})

	handler := rateLimitMiddleware(limiter, mux)

	log.Println("Server started on http://localhost:8080")

	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
