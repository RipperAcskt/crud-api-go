package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/RipperAcskt/crud-api-go/internal/restapi"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	url := "postgres://ripper:150403@localhost:5432/ripper"

	var app restapi.AppHandler
	app.DB = openDB(url)
	defer app.DB.Close()

	mux := http.NewServeMux()
	mux.Handle("/users", http.HandlerFunc(app.Controller))
	log.Fatal(http.ListenAndServe("localhost:8080", mux))

}

func openDB(url string) *sql.DB {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return db
}
