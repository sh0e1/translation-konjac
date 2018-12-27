package middleware

import (
	"net/http"

	"google.golang.org/appengine"
)

// AppEngineContext ...
func AppEngineContext(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(appengine.NewContext(r)))
	}
}
