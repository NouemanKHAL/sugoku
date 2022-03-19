package middleware

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func LogMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}
