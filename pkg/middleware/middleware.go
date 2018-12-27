package middleware

import "net/http"

type (
	middlewareFunc func(http.HandlerFunc) http.HandlerFunc
	middleware     []middlewareFunc
)

// Chain ...
func Chain(mw ...middlewareFunc) middleware {
	return middleware(mw)
}

// Then ...
func (m middleware) Then(h http.HandlerFunc) http.HandlerFunc {
	for i := range m {
		h = m[len(m)-1-i](h)
	}
	return h
}

type contextKey int

// ...
const (
	channelTokenContextKey contextKey = iota + 1
)
