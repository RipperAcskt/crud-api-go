package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/RipperAcskt/crud-api-go/db"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	url := "postgres://ripper:150403@localhost:5432/ripper"

	var db db.Database
	db.DbObject = openDB(url)
	defer db.DbObject.Close()

	mux := http.NewServeMux()
	mux.Handle("/list", http.HandlerFunc(db.List))
	mux.Handle("/create", http.HandlerFunc(db.Create))
	mux.Handle("/delete", http.HandlerFunc(db.Delete))
	mux.Handle("/updateAll", http.HandlerFunc(db.UpdateAll))
	mux.Handle("/update", http.HandlerFunc(db.Update))
	log.Fatal(http.ListenAndServe("localhost:8080", mux))

}

func openDB(url string) *sql.DB {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	return db
}
