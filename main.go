package main

import (
	"fmt"
	"time"

	"github.com/RenanAlmeida225/concurrent-movie-uploader/infra"
	"github.com/RenanAlmeida225/concurrent-movie-uploader/service"
)

func main() {
	start := time.Now()
	db := infra.New()
	r := infra.NewRepo(db)
	s := service.New(r)

	movies := make(chan []*infra.Movie, 2)

	go s.ReadCSV("movies.csv", movies)

	s.SaveMovies(movies)

	elapsed := time.Since(start)
	fmt.Printf("Tempo de execução: %s\n", elapsed)
}
