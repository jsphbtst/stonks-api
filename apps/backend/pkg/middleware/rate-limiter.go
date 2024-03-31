package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const FIVE_SECOND_CTX_TIMEOUT time.Duration = 5 * time.Second
const SIXTY_SECOND_EXPIRY time.Duration = 60 * time.Second

func RateLimitMiddleware(redisClient *redis.Client, limit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cacheCtx, cancel := context.WithTimeout(context.Background(), FIVE_SECOND_CTX_TIMEOUT)
			defer cancel()

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			key := "rate_limit:" + ip
			val, err := redisClient.Get(cacheCtx, key).Result()
			if err == redis.Nil {
				redisClient.Set(cacheCtx, key, 1, SIXTY_SECOND_EXPIRY)
				next.ServeHTTP(w, r)
			} else if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				payload := fmt.Sprintf("{\"message\": \"%s\"}", err.Error())
				w.Write([]byte(payload))
				return
			} else {
				count, _ := strconv.Atoi(val)
				if count >= limit {
					w.WriteHeader(http.StatusBadRequest)
					payload := fmt.Sprintf("{\"message\": \"%s\"}", "Exceeded request limit")
					w.Write([]byte(payload))
					return
				}
				redisClient.Incr(cacheCtx, key)
				next.ServeHTTP(w, r)
			}
		})
	}
}
