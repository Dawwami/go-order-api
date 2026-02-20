# Session Progress — 20 Feb 2026

## Yang Diselesaikan Hari Ini

### Phase 2 — Step 6: Repository Layer (DONE)
- [x] Fix bug `order_repo.go:25` — `Find(orders)` → `Find(&orders)`
- [x] Implement `FindByUserID` di `order_repo.go` — return `([]model.Order, error)` bukan single order
- [x] Implement `user_repo.go` — `Create`, `FindByID`, `FindByEmail`

### Phase 2 — Step 9: AutoMigrate (DONE)
- [x] Wire `db.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{})` di `main.go`
- [x] Import package `model` di `main.go`
- [x] Urutan migrate: `User → Product → Order` (parent sebelum child, karena Order punya FK ke keduanya)

### Phase 2 — Step 7: Service Layer (DONE)
- [x] `internal/service/product_service.go` — `GetAll`, `GetByID`, `Create`, `Update`, `Delete`
- [x] `internal/service/auth_service.go` — `Register` (bcrypt hash), `Login` (verify + JWT), constructor terima `jwtSecret`
- [x] `internal/service/order_service.go` — `Create`, `GetAll`, `GetByID`, `GetByUserID`

### Phase 2 — Step 8: Handler Layer (PARTIAL)
- [x] `internal/handler/product_handler.go` — `Create`, `GetAll`, `GetByID`, `Update`, `Delete`
- [ ] `internal/handler/auth_handler.go` — **BELUM DIBUAT** ← lanjut dari sini
- [ ] `internal/handler/order_handler.go` — belum dibuat

---

## Checkpoint: Lanjut dari Sini

**File berikutnya: `internal/handler/auth_handler.go`**

Struct yang dibutuhkan:
```go
type AuthHandler struct {
    service *service.AuthService
}
```

Dua method yang perlu diimplementasi:
- `Register(c *gin.Context)` — bind JSON (email + password) → `service.Register` → return 201
- `Login(c *gin.Context)` — bind JSON (email + password) → `service.Login` → return token di response data

Setelah `auth_handler.go` selesai, lanjut ke:
1. `internal/handler/order_handler.go`
2. **Phase 2 Step 9** — Route groups di `main.go` (wiring semua handler ke router)
3. **Phase 3** — Custom error (`internal/errors/errors.go`) + recovery middleware

---

## Hal yang Perlu Diingat

### Pola Umum Handler
```go
func (h *XxxHandler) Method(c *gin.Context) {
    var req model.Xxx
    if err := c.ShouldBindJSON(&req); err != nil {
        ErrorResponse(c, http.StatusBadRequest, "bad request")
        return
    }

    ctx := c.Request.Context()
    result, err := h.service.Method(ctx, &req)
    if err != nil {
        ErrorResponse(c, http.StatusInternalServerError, "pesan error")
        return   // ← jangan lupa return!
    }

    SuccessResponse(c, http.StatusOK, result)
}
```

### Parsing ID dari URL param
```go
id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// lalu cast ke uint saat dipakai:
h.service.GetByID(ctx, uint(id))
```

### bcrypt di auth_service.go
- `Register` → `bcrypt.GenerateFromPassword` dulu, baru simpan ke DB
- `Login` → `bcrypt.CompareHashAndPassword([]byte(hashDariDB), []byte(passwordDariUser))`
- **Jangan** generate hash baru di Login — hash sudah ada di DB dari saat Register

---

## Status Keseluruhan Project

| Phase | Step | Status |
|-------|------|--------|
| Phase 1 | Foundation (defer, once, interface{}) | ✅ Done |
| Phase 2 | Step 5: GORM models | ✅ Done |
| Phase 2 | Step 6: Repository layer | ✅ Done |
| Phase 2 | Step 7: Service layer | ✅ Done |
| Phase 2 | Step 8: Handler layer | 🔄 1/3 done (product ✅, auth ❌, order ❌) |
| Phase 2 | Step 9: Router + AutoMigrate | 🔄 AutoMigrate ✅, router wiring ❌ |
| Phase 3 | Error handling + recovery middleware | ❌ Belum |
| Phase 4 | Auth middleware (JWT + context.WithValue) | ❌ Belum |
| Phase 5 | Concurrency (worker, channel, graceful shutdown, mutex) | ❌ Belum |
| Phase 6 | Docker + Testing | ❌ Belum |
