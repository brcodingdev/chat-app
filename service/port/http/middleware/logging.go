package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// ResponseWriter ...
type ResponseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{ResponseWriter: w}
}

// Status ...
func (rw *ResponseWriter) Status() int {
	return rw.status
}

// WriteHeader ...
func (rw *ResponseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

}

// Hijack ...
func (rw *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

// RegisterLoggingMiddleware to HTTP request and duration.
func RegisterLoggingMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					if os.Getenv("LOG_PANIC_TRACE") == "true" {
						level.Error(logger).Log(
							"err", err,
							"trace", debug.Stack(),
						)
					} else {
						level.Error(logger).Log(
							"err", err,
						)
					}
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			level.Info(logger).Log(
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}
