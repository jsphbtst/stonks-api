package services

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jsphbtst/stonks-api/apps/backend/pkg/types"
)

func GetCompanyBySymbol(symbol string) (*types.Companies, error) {
	stmt, err := db.Client.Prepare("SELECT symbol, name, about, sector, industry, mission, vision FROM Companies WHERE symbol = ?")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to prepare statement: %+v\n", err)
		os.Exit(1)
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

	return &company, nil
}