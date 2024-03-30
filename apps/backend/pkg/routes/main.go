package routes

import (
	"database/sql"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"github.com/redis/go-redis/v9"
)

type Services struct {
	SqlClient     *sql.DB
	RedisClient   *redis.Client
	AlgoliaClient *search.Client
	AlgoliaIndex  *search.Index
}

var db = &Services{}

func Init(sqlDb *sql.DB, rc *redis.Client, ac *search.Client, ai *search.Index) {
	db.SqlClient = sqlDb
	db.RedisClient = rc
	db.AlgoliaClient = ac
	db.AlgoliaIndex = ai
}
