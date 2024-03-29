package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/services"
	"github.com/jsphbtst/stonks-api/apps/backend/pkg/types"
)

func GetCompanyById(w http.ResponseWriter, r *http.Request) {
	symbol := chi.URLParam(r, "symbol")

	w.Header().Set("Content-Type", "application/json")

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

	payload := map[string]*types.Companies{"data": company}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		payload := fmt.Sprintf("{\"message\": \"%s\"}", err.Error())
		w.Write([]byte(payload))
		return
	}

	w.Write(jsonData)
}
