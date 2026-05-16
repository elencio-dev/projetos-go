package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limite   int
}

func NewRateLimiter(limite int) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limite:   limite,
	}
}

func (rl *RateLimiter) Permitir(user string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	tempoAgora := time.Now()
	janela := tempoAgora.Add(-1 * time.Minute)

	recents := []time.Time{}

	for _, t := range rl.requests[user] {
		if t.After(janela) {
			recents = append(recents, t)
		}
	}

	rl.requests[user] = recents

	if len(recents) >= rl.limite {
		return false
	}

	rl.requests[user] = append(rl.requests[user], tempoAgora)
	return true

}

func RateLimit(rl *RateLimiter, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, _, _ := r.BasicAuth()
		if !rl.Permitir(user) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}
