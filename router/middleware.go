package router

import (
	"expvar"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/datewu/gtea/handler"
	"golang.org/x/time/rate"
)

// Middleware is a function that takes a http.HandlerFunc and returns a http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

func (ro *bag) enabledCORS(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		// Add the "Vary: Access-Control-Request-Method" header.
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")

		if origin != "" {
			for i := range ro.config.CORS.TrustedOrigins {
				if origin == ro.config.CORS.TrustedOrigins[i] {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// Check if the request has the HTTP method OPTIONS and contains the
					// "Access-Control-Request-Method" header. If it does, then we treat
					// it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response headers, as discussed
						// previously.
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// Write the headers along with a 200 OK status and return from
						// the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}
		next(w, r)
	}
	return middle
}
func (ro *bag) rateLimit(next http.HandlerFunc) http.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		clients = make(map[string]*client)
		mu      sync.Mutex
	)
	delOld := func(interval time.Duration) {
		for {
			time.Sleep(interval)
			mu.Lock()
			for k, v := range clients {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(clients, k)
				}
			}
			mu.Unlock()
		}
	}
	go delOld(time.Minute)

	middle := func(w http.ResponseWriter, r *http.Request) {
		h := handler.NewHandleHelper(w, r)
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			h.ServerErr(err)
			return
		}
		mu.Lock()
		if _, existed := clients[ip]; !existed {
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(ro.config.Limiter.Rps),
					ro.config.Limiter.Burst),
			}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			h.RateLimitExceede()
			return
		}
		mu.Unlock()
		next(w, r)
	}
	return middle
}

func (ro *bag) metrics(next http.HandlerFunc) http.HandlerFunc {
	totalRequestReceived := expvar.NewInt("total_requests_received")
	totalResponsesSend := expvar.NewInt("total_responses_send")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_us")
	middle := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		totalRequestReceived.Add(1)
		next(w, r)
		totalResponsesSend.Add(1)
		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	}
	return middle
}

func recoverPanic(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				handler.WriteJSON(w, http.StatusInternalServerError, handler.Envelope{"recover": err}, nil)
				return
			}
		}()
		next(w, r)
	}
	return middle
}
