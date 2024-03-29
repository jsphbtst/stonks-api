package services

import (
	"context"
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	SqlClient   *sql.DB
	RedisClient *redis.Client
	Ctx         context.Context
}

var db = &Services{}

func Init(sqlDb *sql.DB, rc *redis.Client, ctx context.Context) {
	db.SqlClient = sqlDb
	db.RedisClient = rc
	db.Ctx = ctx
}
