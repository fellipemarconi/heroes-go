package main

import (
	"api/internal/infra/auth"
	dbconn "api/internal/infra/db"
	sqlc "api/internal/infra/db/sqlc"
	redisconn "api/internal/infra/redis"
	"api/internal/server"
	"net/http"
)

func main() {
	pool, err := dbconn.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)
	redisconn.ConnectRedis()
	auth.GenerateJWT()

	r := server.NewRouter(queries)
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
