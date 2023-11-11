package middleware

import (
	"net/http"
	"time"
)

// graceful shutdownのために、requestを遅延させる
func Sleep(t time.Duration, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(t)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
