package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type DB struct {
	Client *sql.DB
}

var db = &DB{}

func Init(uri string) (*sql.DB, error) {
	sqlDb, err := sql.Open("libsql", uri)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", uri, err)
		return nil, err
	}

	db.Client = sqlDb
	return sqlDb, nil
}
