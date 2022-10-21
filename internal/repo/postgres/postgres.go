package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Postgres struct {
	db *sql.DB
}

func New(url string) (*Postgres, error) {
	var p Postgres

	db, err := sql.Open("pgx", url)
	if err != nil {
		return &p, fmt.Errorf("open faild: %v", err)
	}

	p.db = db
	return &p, nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}
