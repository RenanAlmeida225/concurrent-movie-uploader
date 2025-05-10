package infra

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
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

func New() *pgxpool.Pool {
	connStr := "postgres://postgres:postgres@localhost:5432/populate_db"
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Erro ao parsear config: %v", err)
	}

	dbpool, err := pgxpool.New(context.Background(), config.ConnString())
	if err != nil {
		log.Fatalf("Erro ao criar pool: %v", err)
	}

	// Criação da tabela
	_, err = dbpool.Exec(context.Background(), table)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	return dbpool
}
