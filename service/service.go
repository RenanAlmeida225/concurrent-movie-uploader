package service

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/RenanAlmeida225/concurrent-movie-uploader/infra"
)

const MAX_INSERT = 5000

var yearRegex = regexp.MustCompile(`\((\d{4})\)`)

type IService interface {
	ReadCSV(filename string, cMovies chan<- []*infra.Movie)
	SaveMovies(cMovies <-chan []*infra.Movie)
}

type service struct {
	repo infra.IRepository
}

func New(repo infra.IRepository) *service {
	return &service{repo: repo}
}

func (s *service) ReadCSV(filename string, cMovies chan<- []*infra.Movie) {
	start := time.Now()
	var movies []*infra.Movie

	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		log.Fatalf("%s", err)
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%s", err)
		}

		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatalf("%s", err)
		}

		titleYear := s.separateTitleYear(record[1])

		year, err := strconv.Atoi(titleYear[1])
		if err != nil {
			log.Fatalf("%s", err)
		}

		genres := s.separateGenres(record[2])

		movies = append(movies, &infra.Movie{Id: id, Title: titleYear[0], Year: year, Genres: genres})

		if len(movies) == MAX_INSERT {
			cMovies <- movies
			movies = nil
		}
	}

	if len(movies) > 0 {
		cMovies <- movies
	}

	close(cMovies)
	fmt.Printf("Tempo de leitura e envio dos filmes: %s\n", time.Since(start))
}

func (s *service) SaveMovies(cMovies <-chan []*infra.Movie) {
	var wg sync.WaitGroup
	const workers = 5

	worker := func(id int, jobs <-chan []*infra.Movie) {
		defer wg.Done()
		for movies := range jobs {
			start := time.Now()
			if err := s.repo.SaveMultiplesMovies(movies); err != nil {
				log.Printf("worker %d: erro ao salvar filmes: %s", id, err)
			}
			fmt.Printf("Worker %d salvou %d filmes em %s\n", id, len(movies), time.Since(start))
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go worker(i+1, cMovies)
	}

	wg.Wait()

}

func (s *service) separateTitleYear(titleYear string) [2]string {
	matches := yearRegex.FindAllStringSubmatch(titleYear, -1)

	if len(matches) == 0 {
		return [2]string{strings.TrimSpace(titleYear), "0000"}
	}

	// Pega o último match (último ano válido)
	year := matches[len(matches)-1][1]

	// Remove esse último "(YYYY)" do título
	last := strings.LastIndex(titleYear, "("+year+")")
	title := strings.TrimSpace(titleYear[:last])

	return [2]string{title, year}
}

func (s *service) separateGenres(genres string) []string {
	return strings.Split(genres, "|")
}
