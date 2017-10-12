package engine

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Postgres dbEngine
type Postgres struct {
	db *sql.DB
}

// NewPostgresEngine returns engine for Postgres db
func NewPostgresEngine(dsn string) Engine {
	db, err := sql.Open("posgres", dsn)
	if err != nil {
		panic(err)
	}

	return &Postgres{
		db: db,
	}
}

func (p *Postgres) Truncate(table string) error {
	cmd := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)
	fmt.Printf("[Postgres] Executing command %s\n", cmd)

	_, err := p.db.Exec(cmd)
	return err
}
