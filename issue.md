# Setup Initial Golang API Service Structure

## Objective

Membuat struktur awal project Golang API service yang sederhana, rapi, dan scalable untuk kebutuhan pengembangan jangka menengah. Project akan menggunakan pendekatan **feature-based clean architecture** yang tidak terlalu kompleks, dengan pemisahan tanggung jawab yang jelas antara controller, service, repository, model, dan DTO.

Project ini akan menggunakan:
- **Fiber v3** sebagai HTTP framework
- **slog** sebagai structured logger
- **PostgreSQL** sebagai primary database via GORM
- **Redis** sebagai cache / shared state
- **RabbitMQ** sebagai message broker
- **golang-migrate** untuk database migration
- **Swagger** untuk dokumentasi API

## Goals

- Menyusun folder structure awal yang konsisten dan mudah dikembangkan
- Menyiapkan baseline arsitektur project berbasis fitur
- Memisahkan entrypoint **API** dan **worker/consumer** RabbitMQ
- Menambahkan dukungan migration dan Swagger generation via Makefile
- Menetapkan shared infrastructure untuk logger, database, cache, broker, validator, dan response handling

## Required Packages

Package yang wajib digunakan dalam project ini:

```go
github.com/gofiber/fiber/v3
gorm.io/gorm
gorm.io/driver/postgres
github.com/go-playground/validator/v10
github.com/golang-migrate/migrate/v4
github.com/redis/go-redis/v9
github.com/rabbitmq/amqp091-go
log/slog
github.com/google/uuid
```

Tambahan package yang direkomendasikan:

```go
github.com/gofiber/contrib/v3/swaggo
github.com/swaggo/swag/cmd/swag
```

## Proposed Folder Structure

```text
.
├── cmd
│   ├── api
│   │   └── main.go
│   └── worker
│       └── main.go
│
├── internal
│   ├── app
│   │   ├── bootstrap.go
│   │   └── container.go
│   │
│   ├── platform
│   │   ├── config
│   │   │   └── config.go
│   │   ├── database
│   │   │   └── postgres.go
│   │   ├── cache
│   │   │   └── redis.go
│   │   ├── broker
│   │   │   ├── rabbitmq.go
│   │   │   ├── publisher.go
│   │   │   └── consumer.go
│   │   └── logger
│   │       └── slog.go
│   │
│   ├── http
│   │   ├── middleware
│   │   │   ├── recovery.go
│   │   │   ├── request_id.go
│   │   │   ├── logging.go
│   │   │   └── auth.go
│   │   └── routes
│   │       ├── api.go
│   │       ├── health.go
│   │       └── swagger.go
│   │
│   ├── shared
│   │   ├── response
│   │   │   └── response.go
│   │   ├── errors
│   │   │   └── errors.go
│   │   ├── validator
│   │   │   └── validator.go
│   │   ├── pagination
│   │   │   └── pagination.go
│   │   └── utils
│   │       └── utils.go
│   │
│   └── modules
│       └── user
│           ├── controller.go
│           ├── service.go
│           ├── repository.go
│           ├── model.go
│           ├── dto.go
│           └── router.go
│
├── db
│   └── migrations
│       ├── 000001_init.up.sql
│       └── 000001_init.down.sql
│
├── docs
│   └── swagger
│       ├── docs.go
│       ├── swagger.json
│       └── swagger.yaml
│
├── scripts
│   └── wait-for-deps.sh
│
├── .env.example
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

## Penataan `internal/`: batas visibilitas, bukan satu lapisan bisnis

Di Go, `internal/` adalah **batas impor**: kode di dalamnya hanya boleh diimpor oleh paket dalam pohon modul yang sama. Itu **bukan** pernyataan bahwa `app`, `platform`, `http`, dan `modules` adalah jenis concern yang sama—hanya bahwa **semuanya privat untuk service ini**.

- **`internal/app`** — wiring aplikasi: bootstrap, dependency injection, container.
- **`internal/platform`** — konfigurasi dan adapter ke sistem eksternal (PostgreSQL, Redis, RabbitMQ, logger/slog). Bukan fitur bisnis.
- **`internal/http`** — lapisan transport HTTP: middleware dan registrasi route.
- **`internal/shared`** — helper yang dipakai lintas fitur (response, errors, validator, pagination, utils).
- **`internal/modules/<fitur>`** — modul bisnis per domain (contoh: `user`).

### Jika fitur banyak

**Opsi A — tetap datar:** menambah folder fitur langsung di bawah `internal/` masih valid, tetapi root `internal/` cepat ramai.

**Opsi B (disarankan untuk skala besar):** memakai `internal/modules/<fitur>` agar fitur bisnis terkumpul, sementara root `internal/` hanya berisi `app`, `platform`, `http`, `shared`, `modules`.

### Yang sebaiknya dihindari

- Menduplikasi middleware lintas fitur di setiap modul (kecuali benar-benar khusus fitur).
- Struktur **hanya per layer** (`internal/controller`, `internal/service`, …) yang memecah satu fitur ke banyak folder.
- Terlalu banyak folder fitur "longgar" di root `internal/` tanpa namespace `modules/` ketika proyek sudah besar.

### Ringkasan mental model

| Folder | Peran |
|--------|--------|
| `app` | Wiring |
| `platform` | Konfigurasi dan infrastruktur teknis |
| `http` | Middleware dan route HTTP |
| `shared` | Kode umum reusable |
| `modules` | Fitur/domain bisnis |

## Architectural Notes

### 1. Feature-Based Structure
Setiap domain/fitur ditempatkan dalam folder sendiri di dalam `internal/modules/<nama-fitur>/`.

Contoh:

```text
/internal/modules/user
  controller.go
  service.go
  repository.go
  model.go
  dto.go
  router.go
```

Keterangan:
- `controller.go`: menangani HTTP request/response
- `service.go`: berisi business logic
- `repository.go`: akses data ke PostgreSQL
- `model.go`: **DB model only**
- `dto.go`: request/response DTO
- `router.go`: registrasi route per feature

### 2. Model Policy
`model.go` hanya digunakan untuk representasi database model, bukan untuk request/response API dan bukan untuk domain entity terpisah.

### 3. API & Worker Separation
RabbitMQ consumer dipisahkan dari HTTP API dalam entrypoint terpisah:
- `cmd/api/main.go` → HTTP API
- `cmd/worker/main.go` → RabbitMQ consumer/worker

Hal ini mempermudah deployment, scaling, dan lifecycle management.

### 4. Platform (konfigurasi dan infrastruktur teknis)
Semua koneksi dan dependency teknis diletakkan di `internal/platform`, meliputi:
- `config` — load dan parsing konfigurasi
- `database` — PostgreSQL connection
- `cache` — Redis client
- `broker` — RabbitMQ connection / publisher / consumer
- `logger` — slog logger setup

### 5. Shared Utility Layer
Semua helper lintas fitur diletakkan di `internal/shared`, seperti:
- standard API response
- application error mapping
- validator wrapper
- pagination helper
- shared utility functions

### 6. Lapisan HTTP
Middleware dan agregasi route HTTP berada di `internal/http/middleware` dan `internal/http/routes`, terpisah dari modul bisnis di `internal/modules`.

## Swagger Setup

Swagger perlu disiapkan agar dokumentasi API dapat digenerate dari komentar code.

### Recommended Setup
- Swagger generator source: `./cmd/api/main.go`
- Swagger output folder: `./docs/swagger`
- Swagger route: `/swagger/*`

### Notes
- Gunakan generator `swag`
- Gunakan middleware Fiber untuk expose Swagger UI
- Hasil generate Swagger sebaiknya disimpan ke repository agar konsisten untuk tim dan CI

## Makefile Targets

Makefile minimal perlu menyediakan target berikut:

- `run-api` → menjalankan HTTP API
- `run-worker` → menjalankan worker RabbitMQ
- `build-api` → build binary API
- `build-worker` → build binary worker
- `test` → menjalankan seluruh unit test
- `tidy` → merapikan dependency module
- `swagger` → generate Swagger docs
- `migrate-create` → membuat file migration baru
- `migrate-up` → menjalankan migration up
- `migrate-down` → rollback 1 migration
- `migrate-force` → force version migration

## Example Makefile

```makefile
APP_NAME=app
API_MAIN=./cmd/api/main.go
WORKER_MAIN=./cmd/worker/main.go
MIGRATION_DIR=./db/migrations
SWAGGER_OUT=./docs/swagger

.PHONY: run-api
run-api:
	go run $(API_MAIN)

.PHONY: run-worker
run-worker:
	go run $(WORKER_MAIN)

.PHONY: build-api
build-api:
	go build -o bin/api $(API_MAIN)

.PHONY: build-worker
build-worker:
	go build -o bin/worker $(WORKER_MAIN)

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: swagger
swagger:
	swag init -g ./cmd/api/main.go -o $(SWAGGER_OUT) --parseDependency --parseInternal

.PHONY: migrate-create
migrate-create:
	@read -p "migration name: " name; \
	migrate create -ext sql -dir $(MIGRATION_DIR) -seq $$name

.PHONY: migrate-up
migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$$DATABASE_URL" up

.PHONY: migrate-down
migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$$DATABASE_URL" down 1

.PHONY: migrate-force
migrate-force:
	@read -p "version: " version; \
	migrate -path $(MIGRATION_DIR) -database "$$DATABASE_URL" force $$version
```

## Scope of Work

### In Scope
- Menyusun folder structure project
- Menyiapkan baseline dependency wiring
- Menyediakan entrypoint API dan worker
- Menyediakan struktur feature module
- Menyediakan shared infrastructure layer
- Menyiapkan struktur migrations
- Menyiapkan Swagger integration
- Menambahkan Makefile untuk run/build/migrate/swagger

### Out of Scope
- Implementasi business logic feature secara lengkap
- Implementasi authentication/authorization lengkap
- Implementasi observability penuh (metrics/tracing)
- Docker / CI pipeline setup
- Retry/DLQ strategy RabbitMQ yang kompleks
- Caching strategy detail per feature

## Recommendations / Tweaks

1. **Gunakan manual dependency injection**, bukan framework DI, agar wiring tetap jelas dan sederhana.
2. **Tambahkan request ID middleware** sejak awal agar logging dengan `slog` lebih mudah ditelusuri.
3. **Pisahkan DTO dari DB model** secara konsisten.
4. **Gunakan SQL migration sebagai source of truth**, bukan `AutoMigrate` sebagai jalur utama.
5. **Tambahkan health endpoint** minimal `/health` dan `/ready`.
6. **Siapkan graceful shutdown** untuk API maupun worker.
7. **Pisahkan publisher dan consumer RabbitMQ** walaupun masih dalam package broker yang sama.
8. **Commit hasil generate Swagger** agar dokumentasi API tidak bergantung pada local generation saja.

## Missing Points / Technical Considerations

Beberapa hal penting yang perlu diperhatikan saat implementasi:

- Standarisasi format response API
- Standarisasi application error dan HTTP status mapping
- Penamaan migration file yang konsisten
- Penentuan penggunaan Redis sejak awal (cache, lock, idempotency, rate-limit, dsb.)
- Konvensi naming untuk package, file, dan method
- Mekanisme shutdown koneksi PostgreSQL, Redis, dan RabbitMQ
- Strategi idempotency untuk RabbitMQ consumer
- Struktur environment variable dan `.env.example`

## Suggested Acceptance Criteria

- [ ] Project memiliki folder structure sesuai proposal
- [ ] Tersedia `cmd/api/main.go`
- [ ] Tersedia `cmd/worker/main.go`
- [ ] Tersedia minimal satu contoh feature module di `internal/modules/<feature>`
- [ ] Tersedia setup PostgreSQL di `internal/platform/database`
- [ ] Tersedia setup Redis di `internal/platform/cache`
- [ ] Tersedia setup RabbitMQ di `internal/platform/broker`
- [ ] Tersedia setup logger berbasis `slog` di `internal/platform/logger`
- [ ] Tersedia shared validator wrapper
- [ ] Tersedia migration folder dan contoh migration file
- [ ] Tersedia route health check
- [ ] Tersedia Swagger route dan hasil generate docs
- [ ] Tersedia Makefile dengan target run, build, migrate, dan swagger
- [ ] `model.go` digunakan khusus untuk DB model
- [ ] API dan worker dapat dijalankan secara terpisah

## Suggested Task Breakdown

- [ ] Inisialisasi project Go module
- [ ] Tambahkan seluruh dependency utama
- [ ] Buat struktur folder dasar
- [ ] Implementasikan config loader (`internal/platform/config`)
- [ ] Implementasikan slog logger setup
- [ ] Implementasikan PostgreSQL connection setup
- [ ] Implementasikan Redis connection setup
- [ ] Implementasikan RabbitMQ connection setup
- [ ] Implementasikan app bootstrap/container sederhana
- [ ] Implementasikan middleware dasar
- [ ] Implementasikan shared response & error helper
- [ ] Implementasikan validator wrapper
- [ ] Implementasikan contoh feature `user`
- [ ] Implementasikan route aggregator
- [ ] Tambahkan health endpoint
- [ ] Integrasikan Swagger
- [ ] Tambahkan Makefile
- [ ] Tambahkan `.env.example`
- [ ] Tambahkan README setup singkat

## Definition of Done

Issue ini dianggap selesai apabila:
- struktur project sudah terbentuk (termasuk `internal/platform`, `internal/http`, `internal/modules`),
- API entrypoint dan worker entrypoint sudah tersedia,
- Swagger bisa digenerate,
- migration bisa dibuat dan dijalankan via Makefile,
- serta project sudah siap digunakan sebagai baseline pengembangan feature selanjutnya.
