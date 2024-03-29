package services

import "database/sql"

type DB struct {
	Client *sql.DB
}

var db = &DB{}

func Init(sqlDb *sql.DB) {
	db.Client = sqlDb
}
