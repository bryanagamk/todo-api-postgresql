# Project A — Todo API (Mysql)

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
- Docker Compose (app + mysql)
- Makefile untuk perintah rutin