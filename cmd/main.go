package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/RipperAcskt/crud-api-go/internal/config"
	"github.com/RipperAcskt/crud-api-go/internal/repo/postgres"
	"github.com/RipperAcskt/crud-api-go/internal/restapi"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config new faild: %v", err)
	}

	pg, err := postgres.New(cfg.Url)
	if err != nil {
		log.Fatalf("postgres new faild: %v", err)
	}

	app := restapi.New(pg)

	defer app.Close()

	mux := http.NewServeMux()
	mux.Handle("/users", http.HandlerFunc(app.CheckMethod))

	address := cfg.Ip + ":" + cfg.Port
	log.Fatal(http.ListenAndServe(address, mux))
}
