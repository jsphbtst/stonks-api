package middleware

import (
	"log"
	"net/http"
	"time"
)

func RouteRuntimeLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		endTime := time.Now()
		diff := endTime.Sub(startTime)
		log.Printf("Runtime for %s %s is %+v\n", r.Method, r.URL.Path, diff)
	})
}
