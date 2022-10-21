package main

import (
	"log"
	"net/http"

	"github.com/RipperAcskt/crud-api-go/internal/repo/postgres"
	"github.com/RipperAcskt/crud-api-go/internal/restapi"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	url := "postgres://ripper:150403@localhost:5432/ripper"

	pg, err := postgres.New(url)
	if err != nil {
		log.Fatalf("postgres new faild: %v", err)
	}

	app := restapi.New(pg)

	defer app.Close()

	mux := http.NewServeMux()
	mux.Handle("/users", http.HandlerFunc(app.CheckMethod))
	log.Fatal(http.ListenAndServe("localhost:8080", mux))

}
