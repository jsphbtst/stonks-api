package routes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/services"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/types"
	"github.com/redis/go-redis/v9"
)

func GetCompanyById(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	val, err := db.RedisClient.Get(ctx, symbol).Result()
	if err != nil && err != redis.Nil {
		payload := fmt.Sprintf("{\"message\": \"%s\"}", err.Error())
		w.Write([]byte(payload))
		return
	}

	if err != redis.Nil {
		var company types.Companies
		if err := json.Unmarshal([]byte(val), &company); err != nil {
			payload := fmt.Sprintf("{\"message\": \"%s\"}", err.Error())
			w.Write([]byte(payload))
			return
		}

		log.Printf("Found %s in cache.\n", symbol)
		payload := map[string]types.Companies{"data": company}
		jsonData, err := json.Marshal(payload)
		if err != nil {
			payload := fmt.Sprintf("{\"message\": \"%s\"}", err.Error())
			w.Write([]byte(payload))
			return
		}

		w.Write(jsonData)
		return
	}

	company, err := services.GetCompanyBySymbol(symbol)
	if err != nil {
		if err == sql.ErrNoRows {
			payload := fmt.Sprintf("{\"data\": \"%+v\"}", nil)
			w.Write([]byte(payload))
			return
		}

		payload := fmt.Sprintf("{\"message\": \"%s\"}", err.Error())
		w.Write([]byte(payload))
		return
	}

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := json.Marshal(company)
		if err != nil {
			log.Println("Failed to marshall company data: ", err)
			return
		}

		err = db.RedisClient.Set(cacheCtx, symbol, data, 120*time.Second).Err()
		if err != nil {
			log.Println("Failed to Redis SetEX: ", err)
			return
		}

		log.Printf("Concurrently set %s to Redis cache\n", symbol)
	}()

	payload := map[string]*types.Companies{"data": company}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		payload := fmt.Sprintf("{\"message\": \"%s\"}", err.Error())
		w.Write([]byte(payload))
		return
	}

	w.Write(jsonData)
}
