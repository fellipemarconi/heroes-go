package main

import (
	dbconn "api/internal/infra/db"
	sqlc "api/internal/infra/db/sqlc"
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

	r := server.NewRouter(queries)
	err = http.ListenAndServe(":3000", r)
	if err != nil {
		return
	}
}
