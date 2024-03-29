package services

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	SqlClient   *sql.DB
	RedisClient *redis.Client
}

var db = &Services{}

func Init(sqlDb *sql.DB, rc *redis.Client) {
	db.SqlClient = sqlDb
	db.RedisClient = rc
}
