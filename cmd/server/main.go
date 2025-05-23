//go:build integration
// +build integration

package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"squirrel/internal/server"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/squirrel?sslmode=disable"
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("ping: %v", err)
	}
	srv := server.NewServer(server.NewSQLStore(db))
	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv.Handler()))
}
