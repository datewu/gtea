package handler

import (
	"compress/gzip"
	"expvar"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Middleware takes a http.HandlerFunc and returns a http.HandlerFunc.
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

// Append middleware after base middleware.
func Append(base, later Middleware) Middleware {
	if later == nil {
		return base
	}
	if base == nil {
		return later
	}
	md := func(next http.HandlerFunc) http.HandlerFunc {
		next = later(next)
		return base(next)
	}
	return md
}

// Insert middleware befer base middlware
func Insert(base, first Middleware) Middleware {
	if first == nil {
		return base
	}
	if base == nil {
		return first
	}
	md := func(next http.HandlerFunc) http.HandlerFunc {
		next = base(next)
		return first(next)
	}
	return md
}

// Aggregate middlewares.
func Aggregate(mds ...Middleware) Middleware {
	if len(mds) < 1 {
		return nil
	}
	md := mds[0]
	for i := 1; i < len(mds); i++ {
		md = Insert(md, mds[i])
	}
	return md
}

// RecoverPanicMiddleware middleware
func RecoverPanicMiddleware(next http.HandlerFunc) http.HandlerFunc {
	middle := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				ServerErrAny(w, err)
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
			if r.RemoteAddr == "" {
				// for httptest.NewRecorder()
				r.RemoteAddr = "httptest.client:35256"
			}
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				ServerErr(w, err)
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
				RateLimitExceede(w)
				return
			}
			mu.Unlock()
			next(w, r)
		}
		return ratelimit
	}
	return mid
}

// CORSMiddleware unblock trustedOrigins domains
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

// TokenMiddleware GetToken then use check token
func TokenMiddleware(check func(string) (bool, error)) Middleware {
	mid := func(next http.HandlerFunc) http.HandlerFunc {
		checkToken := func(w http.ResponseWriter, r *http.Request) {
			token, err := GetToken(r, "token")
			if err != nil {
				if err == ErrNoToken {
					AuthenticationRequire(w)
					return
				}
				InvalidCredentials(w)
				return
			}
			ok, err := check(token)
			if err != nil {
				ServerErr(w, err)
				return
			}
			if !ok {
				NotPermitted(w)
				return
			}
			next(w, r)
		}
		return checkToken
	}
	return mid
}

var gzPool = sync.Pool{
	New: func() any {
		w := gzip.NewWriter(io.Discard)
		gzip.NewWriterLevel(w, gzip.BestCompression)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// GipMiddleware good for serving static html/js/css file
func GzipMiddleware(next http.HandlerFunc) http.HandlerFunc {
	gzipFunc := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next(w, r)
			return
		}
		w.Header().Set("Content-Encoding", "gzip")
		gz := gzPool.Get().(*gzip.Writer)
		defer gzPool.Put(gz)

		gz.Reset(w)
		defer gz.Close()
		next(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	}
	return gzipFunc
}
