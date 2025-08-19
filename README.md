# Project A — Todo API (Postgresql)

Fokus: REST API idiomatik, context, testing, Swagger, pprof, Docker.

### Target Fitur
- CRUD Todo (/api/v1/todos)
- Filtering + pagination (by status, due_date)
- Auth sederhana: JWT (register/login)
- OpenAPI spec + Swagger UI
- Unit & integration test
- pprof untuk profiling & tracing sederhana
  
### Arsitektur & Tooling
- Fiber (router & middleware)
- GORM + pgx (driver) + golang-migrate (migrations)
- OpenAPI pakai swaggo/swag + fiber-swagger
- pprof via net/http/pprof di port terpisah
- Docker Compose (app + Postgresql)
- Makefile untuk perintah rutin

# ToDo API — Day 1 (Go 1.21)

## Goal
- Setup proyek Go (Fiber)
- Entity `Todo` + validasi (tanpa DB)
- Unit test table-driven
- Server minimal + `/health`

## Run
```bash
make run
# curl localhost:8080/health
