package services

import (
	"fmt"
	"os"

	"github.com/jsphbtst/stonks-api/apps/backend/pkg/types"
)

func GetCompanyBySymbol(symbol string) (*types.Companies, error) {
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

	// can be sql.ErrNoRows
	if err != nil {
		return nil, err
	}

	return &company, nil
}
