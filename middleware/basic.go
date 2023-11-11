package middleware

import (
	"net/http"

	"github.com/TechBowl-japan/go-stations/env"
)

// curl -u admin:admin http://localhost:8080/healthz
func BasicAuth(env *env.Env, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		id, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		if id != env.BasicAuthId || password != env.BasicAuthPass {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
