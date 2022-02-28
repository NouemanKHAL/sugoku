package middleware

import "net/http"

type Middleware func(h http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) == 0 {
		return f
	}
	return m[0](Chain(f, m[1:]...))
}
