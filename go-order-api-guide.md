# Go Order Processing API — Learning Guide

> **Tujuan**: Belajar semua konsep Go dari `go-concepts.md` secara hands-on lewat 1 project.
> **Mode**: AI hanya boleh **guide**, TIDAK boleh eksekusi/write code. User coding sendiri.
> **Referensi konsep**: `go-concepts.md` (defer, panic/recover, mutex, goroutine, channel, context, select, waitgroup, once, empty interface, type assertion)

---

## Project Overview

REST API backend untuk **sistem pemesanan** (mirip food delivery) dengan **background worker** yang memproses order secara async.

### Tech Stack

| Tech | Fungsi |
|------|--------|
| **Go** | Language |
| **Gin** | HTTP framework |
| **GORM** | ORM (explicitly required) |
| **PostgreSQL** | Database |
| **JWT** (`golang-jwt/jwt/v5`) | Authentication |
| **bcrypt** (`golang.org/x/crypto/bcrypt`) | Password hashing |
| **Docker + Docker Compose** | Containerization |

### Architecture

Clean Architecture: **Handler → Service → Repository**

---

## Concept Mapping

Setiap konsep Go harus muncul **secara natural** di project, bukan dipaksakan.

| # | Konsep | Lokasi di Project | Penjelasan |
|---|--------|-------------------|------------|
| 1 | **Defer** | `database.go` (close DB), `ratelimiter.go` (unlock mutex), `recovery.go` | Cleanup resource pasti jalan |
| 2 | **Panic & Recover** | `middleware/recovery.go` | Recovery middleware catch panic, return 500 |
| 3 | **Mutex** | `middleware/ratelimiter.go` | Protect `map[string]int` counter per IP dari concurrent access |
| 4 | **Goroutine** | `worker/order_worker.go`, rate limiter reset | Background worker pool, periodic task |
| 5 | **Channel** | `worker/order_worker.go` | `orderQueue chan uint` — handler kirim order ID, worker consume |
| 6 | **Context** | Semua layer (handler→service→repo→GORM) | `c.Request.Context()` → pass ke service → `db.WithContext(ctx)` |
| 7 | **Select** | `main.go` (graceful shutdown), `order_worker.go` (timeout) | Listen OS signal + channel, worker timeout |
| 8 | **WaitGroup** | `main.go`, `order_worker.go` | Tunggu semua worker goroutine selesai sebelum server shutdown |
| 9 | **Once** | `database/database.go` | Singleton GORM `*gorm.DB` instance via `sync.Once` |
| 10 | **Empty Interface** | API response struct `Data interface{}` | Generic response wrapper |
| 11 | **Type Assertion** | Handler error handling | `appErr, ok := err.(*errors.AppError)` untuk custom HTTP status |

---

## Folder Structure

```
go-order-api/
├── cmd/
│   └── server/
│       └── main.go                  # entry point, router, graceful shutdown
├── internal/
│   ├── config/
│   │   └── config.go                # env vars loader
│   ├── model/
│   │   ├── user.go                  # GORM model: id, email, password
│   │   ├── product.go               # GORM model: id, name, price, stock
│   │   └── order.go                 # GORM model: id, user_id, product_id, qty, status
│   ├── handler/
│   │   ├── auth_handler.go          # register, login (JWT)
│   │   ├── product_handler.go       # CRUD products
│   │   └── order_handler.go         # create order → push ke channel
│   ├── service/
│   │   ├── auth_service.go          # hash password, generate JWT
│   │   ├── product_service.go       # business logic product
│   │   └── order_service.go         # business logic order, context passing
│   ├── repository/
│   │   ├── user_repo.go             # GORM queries + ctx
│   │   ├── product_repo.go          # GORM queries + ctx
│   │   └── order_repo.go            # GORM queries + ctx
│   ├── middleware/
│   │   ├── recovery.go              # defer + recover() → 500
│   │   ├── auth.go                  # JWT validation, context.WithValue(userID)
│   │   └── ratelimiter.go           # sync.Mutex + map counter per IP
│   ├── worker/
│   │   └── order_worker.go          # goroutine pool, channel consumer, select, waitgroup
│   ├── errors/
│   │   └── errors.go                # AppError struct, predefined errors
│   └── database/
│       └── database.go              # GORM init via sync.Once, defer close
├── docker-compose.yml               # api + postgres services
├── Dockerfile                       # multi-stage build
├── go.mod
└── go.sum
```

---

## Learning Phases

### Phase 1 — Foundation
**Konsep: Defer, Once, Empty Interface**

| Step | Task | File | Detail |
|------|------|------|--------|
| 1 | `go mod init go-order-api`, install deps | `go.mod` | gin, gorm, gorm/driver/postgres, golang-jwt/jwt/v5, x/crypto/bcrypt |
| 2 | Database singleton | `database/database.go` | `sync.Once` untuk init `*gorm.DB`, function `GetDB()` return singleton |
| 3 | Config loader | `config/config.go` | Load env vars (DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME, JWT_SECRET) |
| 4 | Generic API response | (di handler atau pkg terpisah) | Struct `Response{Success bool, Message string, Data interface{}}`, helper `Success(data)` dan `Error(msg)` |

### Phase 2 — GORM & Clean Architecture
**Konsep: Context**

| Step | Task | File | Detail |
|------|------|------|--------|
| 5 | Define GORM models | `model/*.go` | User, Product, Order dengan GORM tags |
| 6 | Repository layer | `repository/*.go` | Semua method terima `ctx context.Context`, pakai `db.WithContext(ctx)` |
| 7 | Service layer | `service/*.go` | Business logic, terima ctx dari handler, pass ke repo |
| 8 | Handler layer | `handler/*.go` | `ctx := c.Request.Context()`, pass ke service, return JSON |
| 9 | Router + auto migrate | `main.go` | Gin router, group routes `/api/v1`, `db.AutoMigrate(models...)` |

### Phase 3 — Error Handling
**Konsep: Type Assertion, Panic & Recover**

| Step | Task | File | Detail |
|------|------|------|--------|
| 10 | Custom AppError | `errors/errors.go` | Struct `AppError{Code, Message}`, implement `Error() string`, predefined: ErrNotFound(404), ErrUnauthorized(401), ErrBadRequest(400), ErrInternal(500) |
| 11 | Type assertion di handler | `handler/*.go` | `appErr, ok := err.(*errors.AppError)` → return appropriate HTTP status |
| 12 | Recovery middleware | `middleware/recovery.go` | `defer func() { if r := recover(); r != nil { ... c.JSON(500, ...) } }()` |

### Phase 4 — Authentication
**Konsep: Context.WithValue**

| Step | Task | File | Detail |
|------|------|------|--------|
| 13 | Auth handler | `handler/auth_handler.go` | Register: bcrypt hash → save via GORM. Login: verify → generate JWT |
| 14 | Auth middleware | `middleware/auth.go` | Parse `Authorization: Bearer <token>`, extract userID dari claims, `context.WithValue(ctx, userIDKey, userID)`, handler ambil via `ctx.Value(userIDKey)` |

### Phase 5 — Concurrency (INTI PROJECT)
**Konsep: Goroutine, Channel, Select, WaitGroup, Mutex**

| Step | Task | File | Detail |
|------|------|------|--------|
| 15 | Order worker pool | `worker/order_worker.go` | Buffered channel `orderQueue chan uint`, spawn N goroutine, `for orderID := range orderQueue`, process: update status pending→processing→completed via GORM, `context.WithTimeout` per job |
| 16 | Handler → Worker | `handler/order_handler.go` | `POST /orders` → save DB (pending) → `orderQueue <- order.ID` → return 202 Accepted |
| 17 | Graceful shutdown | `main.go` | `quit := make(chan os.Signal, 1)`, `signal.Notify(quit, SIGINT, SIGTERM)`, **select** tunggu quit, close channel, **wg.Wait()** tunggu workers, shutdown server dengan **context.WithTimeout**, **defer** close DB |
| 18 | Rate limiter | `middleware/ratelimiter.go` | `map[string]int` + `sync.Mutex`, lock/unlock per request, goroutine reset counter tiap 1 menit |

### Phase 6 — Polish

| Step | Task | File | Detail |
|------|------|------|--------|
| 19 | Dockerize | `Dockerfile`, `docker-compose.yml` | Multi-stage build, compose: api + postgres |
| 20 | Testing | `*_test.go` | Minimal: test service layer dengan mock repo |

---

## End-to-End Request Flow

```
Client: POST /orders {product_id: 1, quantity: 2}
  │
  ├─► [Recovery Middleware]     defer + recover() ──── panic protection
  ├─► [Rate Limiter]           mutex.Lock() → check map[ip]counter → mutex.Unlock()
  ├─► [Auth Middleware]        JWT verify → context.WithValue(userID)
  │
  ▼
[Order Handler]
  │  ctx := c.Request.Context()
  │  userID := ctx.Value(userIDKey)       ← dari auth middleware
  │  order, err := service.Create(ctx, req)
  │  if err → type assertion: *AppError?  ← custom error handling
  │  orderQueue <- order.ID               ← kirim ke CHANNEL (buffered)
  │  return 202 Accepted
  │
  ▼
[Order Worker Pool] ── N GOROUTINE, range dari channel
  │  for orderID := range orderQueue {
  │      ctx, cancel := context.WithTimeout(bg, 30s)
  │      defer cancel()
  │
  │      select {
  │      case <-ctx.Done():
  │          // timeout, mark failed
  │      default:
  │          // update status → processing
  │          // simulate work (sleep)
  │          // update status → completed
  │          // semua GORM query pakai db.WithContext(ctx)
  │      }
  │  }
  │  wg.Done()   ← saat channel closed & loop selesai
  │
  ▼
[Graceful Shutdown] ── di main.go
  select {
  case <-quit:             ← OS signal (SIGINT/SIGTERM)
      close(orderQueue)    ← trigger worker loop exit
      wg.Wait()            ← tunggu semua worker selesai
      server.Shutdown(ctx) ← context.WithTimeout(5s)
      // defer db.Close()  ← cleanup DB
  }
```

---

## API Endpoints

| Method | Path | Auth? | Handler | Description |
|--------|------|-------|---------|-------------|
| POST | `/api/v1/register` | No | auth | Register user baru |
| POST | `/api/v1/login` | No | auth | Login, return JWT |
| GET | `/api/v1/products` | No | product | List semua product |
| GET | `/api/v1/products/:id` | No | product | Get product by ID |
| POST | `/api/v1/products` | Yes | product | Create product (admin) |
| PUT | `/api/v1/products/:id` | Yes | product | Update product |
| DELETE | `/api/v1/products/:id` | Yes | product | Delete product |
| POST | `/api/v1/orders` | Yes | order | Create order → async process |
| GET | `/api/v1/orders` | Yes | order | List user's orders |
| GET | `/api/v1/orders/:id` | Yes | order | Get order status |

---

## GORM Patterns yang Harus Dipelajari

```
Wajib pakai di project ini:

1. AutoMigrate          → db.AutoMigrate(&User{}, &Product{}, &Order{})
2. Create               → db.WithContext(ctx).Create(&order)
3. First / Find         → db.WithContext(ctx).First(&product, id)
4. Save / Updates       → db.WithContext(ctx).Save(&order)
5. Delete               → db.WithContext(ctx).Delete(&product, id)
6. Where                → db.WithContext(ctx).Where("user_id = ?", uid).Find(&orders)
7. Preload (relations)  → db.WithContext(ctx).Preload("Product").Find(&orders)
8. Transaction          → db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { ... })
9. Hooks (optional)     → BeforeCreate, AfterUpdate
```

---

## Order Status State Machine

```
pending ──► processing ──► completed
                │
                └──► failed (timeout / error)
```

---

## Progress Tracker

Checklist untuk track progress belajar:

- [ ] Phase 1: Foundation (defer, once, empty interface)
  - [ ] Step 1: Project init + deps
  - [ ] Step 2: Database singleton (sync.Once)
  - [ ] Step 3: Config loader
  - [ ] Step 4: Generic API response (interface{})
- [ ] Phase 2: GORM & Clean Architecture (context)
  - [ ] Step 5: GORM models
  - [ ] Step 6: Repository layer (ctx + db.WithContext)
  - [ ] Step 7: Service layer
  - [ ] Step 8: Handler layer
  - [ ] Step 9: Router + auto migrate
- [ ] Phase 3: Error Handling (type assertion, panic/recover)
  - [ ] Step 10: Custom AppError
  - [ ] Step 11: Type assertion di handler
  - [ ] Step 12: Recovery middleware
- [ ] Phase 4: Authentication (context.WithValue)
  - [ ] Step 13: Auth handler (register/login + JWT)
  - [ ] Step 14: Auth middleware (JWT verify + ctx.WithValue)
- [ ] Phase 5: Concurrency (goroutine, channel, select, waitgroup, mutex)
  - [ ] Step 15: Order worker pool
  - [ ] Step 16: Handler → Worker (channel push)
  - [ ] Step 17: Graceful shutdown
  - [ ] Step 18: Rate limiter (mutex)
- [ ] Phase 6: Polish
  - [ ] Step 19: Docker + Docker Compose
  - [ ] Step 20: Testing

---

## Rules for AI Assistant

1. **JANGAN** tulis/execute code. Hanya guide dan jelaskan.
2. **BOLEH** kasih snippet kecil untuk illustrasi konsep, tapi user harus ketik sendiri.
3. **BOLEH** review code yang user tunjukkan dan kasih feedback.
4. Kalau user stuck, kasih hint bertahap — jangan langsung kasih jawaban.
5. Encourage user untuk intentionally bikin error (race condition tanpa mutex, goroutine leak tanpa close channel) supaya paham kenapa konsep itu penting.
