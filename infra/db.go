package infra

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

const table = `
		CREATE TABLE IF NOT EXISTS movies (
			id     INTEGER PRIMARY KEY,
			title  TEXT NOT NULL,
			year   INTEGER NOT NULL,
			genres TEXT[]
		);`

type Movie struct {
	Id     int      `db:"id"`
	Title  string   `db:"title"`
	Year   int      `db:"year"`
	Genres []string `db:"genres"`
}

func New() *sql.DB {
	connStr := "postgres://postgres:postgres@localhost:5432/populate_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Erro ao abrir conexão com o banco: %s", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Erro ao testar conexão com o banco: %s", err)
	}

	if _, err = db.Exec(table); err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	return db
}
