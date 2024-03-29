package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jsphbtst/stonks-api/apps/backend/pkg/types"
	"github.com/redis/go-redis/v9"
)

func GetCompanyBySymbol(symbol string) (*types.Companies, error) {
	val, err := db.RedisClient.Get(db.Ctx, symbol).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	if err != redis.Nil {
		var company types.Companies
		if err := json.Unmarshal([]byte(val), &company); err != nil {
			return nil, err
		}

		log.Printf("Found %s in cache.\n", symbol)
		return &company, nil
	}

	stmt, err := db.SqlClient.Prepare("SELECT symbol, name, about, sector, industry, mission, vision FROM Companies WHERE symbol = ?")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to prepare statement: %+v\n", err)
		return nil, err
	}

	var company types.Companies
	err = stmt.QueryRow(symbol).Scan(
		&company.Symbol,
		&company.Name,
		&company.About,
		&company.Sector,
		&company.Industry,
		&company.Mission,
		&company.Vision,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	data, err := json.Marshal(company)
	if err != nil {
		return nil, err
	}

	err = db.RedisClient.Set(db.Ctx, symbol, data, 120*time.Second).Err()
	if err != nil {
		return nil, err
	}

	return &company, nil
}
