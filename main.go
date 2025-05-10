package main

import (
	"fmt"
	"log"
	"time"

	"os"
	"runtime/pprof"

	"github.com/RenanAlmeida225/concurrent-movie-uploader/infra"
	"github.com/RenanAlmeida225/concurrent-movie-uploader/service"
)

func main() {
	db := infra.New()
	r := infra.NewRepo(db)
	s := service.New(r)

	movies := make(chan []*infra.Movie, 20)

	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	start := time.Now()

	go s.ReadCSV("movies.csv", movies)

	s.SaveMovies(movies)

	elapsed := time.Since(start)
	fmt.Printf("Tempo de execução: %s\n", elapsed)
}
