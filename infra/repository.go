package infra

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IRepository interface {
	SaveMultiplesMovies(movies []*Movie) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) SaveMultiplesMovies(movies []*Movie) error {
	rows := make([][]interface{}, len(movies))
	for i, movie := range movies {
		rows[i] = []interface{}{movie.Id, movie.Title, movie.Year, movie.Genres}
	}

	_, err := r.db.CopyFrom(
		context.Background(),
		pgx.Identifier{"movies"},
		[]string{"id", "title", "year", "genres"},
		pgx.CopyFromRows(rows),
	)

	return err
}
