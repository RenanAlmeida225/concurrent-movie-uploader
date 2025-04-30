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

	"github.com/RenanAlmeida225/concurrent-movie-uploader/infra"
)

const MAX_INSERT = 2500

type IService interface {
	ReadCSV(filename string, cMovies chan<- []infra.Movie)
	separateTitleYear(titleYear string) [2]string
	separateGenres(genres string) []string
	SaveMovies(cMovies <-chan []infra.IRepository)
}

type service struct {
	repo infra.IRepository
}

func New(repo infra.IRepository) *service {
	return &service{repo: repo}
}

func (s *service) ReadCSV(filename string, cMovies chan<- []infra.Movie) {
	var movies []infra.Movie

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

		movies = append(movies, infra.Movie{Id: id, Title: titleYear[0], Year: year, Genres: genres})

		if len(movies) == MAX_INSERT {
			cMovies <- movies
			movies = nil
		}
	}

	if len(movies) > 0 {
		cMovies <- movies
	}

	close(cMovies)
}

// ^(.*)\((\d{4})\)$
func (s *service) separateTitleYear(titleYear string) [2]string {
	re := regexp.MustCompile(`\((\d{4})\)`)
	matches := re.FindAllStringSubmatch(titleYear, -1)

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

func (s *service) SaveMovies(cMovies <-chan []infra.Movie) {
	count := 0
	for movies := range cMovies {
		if err := s.repo.SaveMultiplesMovies(movies); err != nil {
			log.Fatalf("erro ao salvar filmes: %s", err)
		}
		count += len(movies)
		fmt.Printf("Salvo total de %d filmes até agora\n", count)
	}
}
