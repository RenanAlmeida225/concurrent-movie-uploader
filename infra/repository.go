package infra

import (
	"database/sql"
	"log"

	"github.com/lib/pq"
)

type IRepository interface {
	SaveMultiplesMovies(movies []Movie) error
}

type repository struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) SaveMultiplesMovies(movies []Movie) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO movies (id, title, year, genres) VALUES ($1, $2, $3, $4)")
	if err != nil {
		log.Fatalf("erro no filme: %s", "movie.Title")
		return err
	}
	defer stmt.Close()

	for _, movie := range movies {
		if _, err = stmt.Exec(movie.Id, movie.Title, movie.Year, pq.Array(movie.Genres)); err != nil {
			log.Fatalf("erro no filme: %s", movie.Title)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
