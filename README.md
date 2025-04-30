# go-movie-loader 🎬

Um CLI em Go para importar filmes de um arquivo CSV para um banco de dados PostgreSQL com alta performance, usando concorrência com goroutines.

## 📌 Visão Geral

Este projeto foi desenvolvido com foco em desempenho e escalabilidade. Utiliza goroutines e canais para leitura e processamento concorrente de grandes volumes de dados contidos no CSV de filmes, realizando a persistência eficiente no banco de dados PostgreSQL.

## 🗂️ Estrutura do Projeto

. ├── app/ # (Reservado para futuros handlers ou lógica extra) ├── docker-compose.yaml # Banco PostgreSQL em container ├── DOCKERFILE # Dockerfile do app Go ├── go.mod / go.sum # Gerenciamento de dependências ├── infra/ │ ├── db.go # Conexão com o banco │ └── repository.go # Inserções e queries ├── main.go # Entry point do CLI ├── movies.csv # Arquivo CSV com dados de filmes └── service/ └── service.go # Regras de negócio e lógica concorrente

## 🧾 Formato Esperado do CSV

movieId,title,genres 1,Toy Story (1995),Adventure|Animation|Children|Comedy|Fantasy 2,Jumanji (1995),Adventure|Children|Fantasy

**A estrutura processada inclui:**

- `id`: ID do filme (ex: 1)
- `title`: Nome do filme (ex: Jumanji)
- `year`: Ano extraído do título (ex: 1995)
- `genres`: Lista de gêneros separados por `|` (ex: Adventure, Fantasy)

## 🚀 Como Executar

### 1. Subir o banco com Docker

```bash
docker-compose up -d

2. Rodar o CLI

go run main.go

Ou, se quiser buildar com Docker:

docker build -t go-movie-loader .
docker run --env-file .env --network host go-movie-loader

⚙️ Configuração

Crie um arquivo .env (ou exporte no terminal) com:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=moviesdb
CSV_FILE_PATH=./movies.csv
WORKERS=10

🧠 Tecnologias e Conceitos

    Go

    Goroutines e Channels

    PostgreSQL

    Docker

    Clean Architecture (divisão em infra, service e main)

    Leitura eficiente de CSV com encoding/csv
```
