package handler

import (
	"expvar"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Middleware is a function that takes a http.HandlerFunc and returns a http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// VoidHandlerFunc ...
var VoidHandlerFunc = http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
})

// VoidMiddleware ...
var VoidMiddleware = Middleware(func(next http.HandlerFunc) http.HandlerFunc {
	return next
})

// AbortMiddleware ...
var AbortMiddleware = Middleware(func(next http.HandlerFunc) http.HandlerFunc {
	return VoidHandlerFunc
})

// AggregateMds aggrate middleware to one
// lower index ms[i] more outer middleware
func AggregateMds(ms []Middleware) Middleware {
	size := len(ms)
	if size == 0 {
		return nil
	}
	md := func(next http.HandlerFunc) http.HandlerFunc {
		h := func(w http.ResponseWriter, r *http.Request) {
			for i := size - 1; i >= 0; i-- {
				if ms[i] == nil {
					continue
				}
				next = ms[i](next)
			}
			next(w, r)
		}
		return h
	}
	return md
}

// RecoverPanicMiddleware middleware
func RecoverPanicMiddleware(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				WriteJSON(w, http.StatusInternalServerError, Envelope{"recover": err}, nil)
				return
			}
		}()
		next(w, r)
	}
	return middle
}

var (
	totalRequestReceived            = expvar.NewInt("total_requests_received")
	totalResponsesSend              = expvar.NewInt("total_responses_send")
	totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_us")
)

// MetricsMiddleware middleware enable expvar profile
func MetricsMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

// RateLimitMiddleware return a middle
func RateLimitMiddleware(rps float64, brust int) Middleware {
	mid := func(next http.HandlerFunc) http.HandlerFunc {
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

		ratelimit := func(w http.ResponseWriter, r *http.Request) {
			h := NewHandleHelper(w, r)
			if r.RemoteAddr == "" {
				// for httptest.NewRecorder()
				r.RemoteAddr = "httptest.client:35256"
			}
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				h.ServerErr(err)
				return
			}
			mu.Lock()
			if _, existed := clients[ip]; !existed {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(rps), brust),
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
		return ratelimit
	}
	return mid
}

func CORSMiddleware(trustedOrigins []string) Middleware {
	mid := func(next http.HandlerFunc) http.HandlerFunc {
		cors := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")

			// Add the "Vary: Access-Control-Request-Method" header.
			w.Header().Add("Vary", "Access-Control-Request-Method")

			origin := r.Header.Get("Origin")

			if origin != "" {
				for i := range trustedOrigins {
					if origin == trustedOrigins[i] {
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
		return cors
	}
	return mid
}

func TokenMiddleware(check func(string) (bool, error)) Middleware {
	mid := func(next http.HandlerFunc) http.HandlerFunc {
		checkToken := func(w http.ResponseWriter, r *http.Request) {
			h := NewHandleHelper(w, r)
			token, err := GetToken(r, "token")
			if err != nil {
				if err == ErrNoToken {
					h.AuthenticationRequire()
					return
				}
				h.InvalidCredentials()
				return
			}
			ok, err := check(token)
			if err != nil {
				h.ServerErr(err)
				return
			}
			if !ok {
				h.NotPermitted()
				return
			}
			next(w, r)
		}
		return checkToken
	}
	return mid
}
