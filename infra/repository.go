package infra

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

type IRepository interface {
	SaveMultiplesMovies(movies []*Movie) error
}

type repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) SaveMultiplesMovies(movies []*Movie) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var (
		valueStrings []string
		valueArgs    []any
	)

	for i, movie := range movies {
		// Cada grupo de valores precisa de placeholders únicos: ($1, $2, $3, $4), ...
		n := i * 4
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", n+1, n+2, n+3, n+4))
		valueArgs = append(valueArgs, movie.Id, movie.Title, movie.Year, pq.Array(movie.Genres))
	}

	stmt := fmt.Sprintf("INSERT INTO movies (id, title, year, genres) VALUES %s", strings.Join(valueStrings, ","))

	if _, err := tx.Exec(stmt, valueArgs...); err != nil {
		log.Printf("erro no bulk insert: %s", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
