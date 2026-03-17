# Go To-Do API

REST API untuk mengelola daftar tugas (to-do list) yang dibangun dengan **Go** dan **PostgreSQL** (Local atau Supabase).

## ✨ Fitur

- **CRUD Lengkap** - Create, Read, Update, Delete tugas
- **Pagination** - Navigasi data dengan pagination
- **Filtering** - Filter tugas berdasarkan status (`pending`/`completed`)
- **Search** - Pencarian tugas berdasarkan judul atau deskripsi
- **Validasi Input** - Validasi semua request body
- **Error Handling** - Penanganan error yang konsisten dengan format JSON
- **Concurrent Execution** - Penggunaan goroutines, WaitGroup, channels, dan mutex
- **Batch Create** - Membuat banyak tugas sekaligus secara concurrent (`POST /tasks/batch`)
- **Structured Logging** - Logging terstruktur menggunakan zerolog
- **Migration System** - Sistem migrasi database seperti Laravel (migrate, rollback, reset, fresh, status)
- **Dual Database** - Mendukung PostgreSQL lokal dan Supabase
- **Graceful Shutdown** - Server berhenti dengan aman saat menerima sinyal shutdown
- **CORS & Recovery** - Middleware untuk CORS dan panic recovery

## 🏗️ Arsitektur

```
backend/
├── main.go                          # Entry point + migration CLI commands
├── config/
│   └── config.go                    # Konfigurasi (dual mode: local / supabase)
├── db/
│   ├── database.go                  # Koneksi database
│   ├── migrator.go                  # Migration engine (seperti Laravel)
│   └── migrations/                  # File-file migrasi SQL
│       ├── 000001_create_tasks_table.up.sql
│       └── 000001_create_tasks_table.down.sql
├── handler/
│   └── task_handler.go              # HTTP handler (controller)
├── middleware/
│   └── middleware.go                # Logger, CORS, Recovery middleware
├── model/
│   └── task.go                      # Model & validasi
├── repository/
│   └── task_repository.go           # Akses database (repository pattern)
├── service/
│   └── task_service.go              # Business logic & concurrent operations
├── go.mod
├── go.sum
├── .env.example
└── README.md
```

## 📋 Prasyarat

1. **Go** versi 1.21 atau lebih baru - [Download Go](https://golang.org/dl/)
2. **PostgreSQL** (pilih salah satu):
   - **Lokal**: PostgreSQL 12+ terinstall di komputer
   - **Supabase**: Akun gratis di [supabase.com](https://supabase.com/)

## 🚀 Cara Menjalankan

### 1. Masuk ke Folder Backend

```bash
cd backend
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Konfigurasi Database (Pilih Salah Satu)

Salin file `.env.example`:

```bash
cp .env.example .env
```

#### Opsi A: PostgreSQL Lokal

1. Pastikan PostgreSQL terinstall dan berjalan
2. Buat database:

```bash
psql -U postgres -c "CREATE DATABASE todo_db;"
```

3. Edit `.env`:

```env
SERVER_PORT=8080

# Kosongkan atau comment DATABASE_URL
# DATABASE_URL=

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=todo_db
DB_SSLMODE=disable
```

#### Opsi B: Supabase

1. Buka [supabase.com/dashboard](https://supabase.com/dashboard) → **New Project**
2. Set nama project dan **database password**
3. Setelah project ready, buka **Settings** → **Database**
4. Copy **Connection string (URI)**
5. Edit `.env`:

```env
SERVER_PORT=8080

# Isi DATABASE_URL (opsi ini akan diprioritaskan)
DATABASE_URL=postgresql://postgres.abcdefghij:YourPassword@aws-0-ap-southeast-1.pooler.supabase.com:6543/postgres
```

> ⚠️ Jika `DATABASE_URL` diisi, maka konfigurasi `DB_HOST`, `DB_PORT`, dll akan **diabaikan**.

### 4. Jalankan Migrasi Database

```bash
go run main.go migrate
```

Output yang diharapkan:

```
📦 Running migrations...
INF Migrations table ready
INF Migrating... migration=000001_create_tasks_table batch=1
INF Migrated successfully ✓ migration=000001_create_tasks_table
INF All migrations completed batch=1 count=1
```

### 5. Jalankan Server

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`.

> 💡 Server juga otomatis menjalankan migrasi pending saat startup.

## 🗃️ Migration Commands

Sistem migrasi ini bekerja seperti Laravel Artisan:

| Command                           | Setara Laravel                 | Deskripsi                                    |
| --------------------------------- | ------------------------------ | -------------------------------------------- |
| `go run main.go migrate`          | `php artisan migrate`          | Jalankan semua migrasi yang belum dijalankan |
| `go run main.go migrate:rollback` | `php artisan migrate:rollback` | Rollback batch migrasi terakhir              |
| `go run main.go migrate:reset`    | `php artisan migrate:reset`    | Rollback semua migrasi                       |
| `go run main.go migrate:fresh`    | `php artisan migrate:fresh`    | Drop semua tabel & jalankan ulang migrasi    |
| `go run main.go migrate:status`   | `php artisan migrate:status`   | Lihat status setiap migrasi                  |

### Contoh Output `migrate:status`

```
┌─────────────────────────────────────────────┬──────────┐
│ Migration                                   │ Status   │
├─────────────────────────────────────────────┼──────────┤
│ 000001_create_tasks_table                   │ ✅ Ran   │
└─────────────────────────────────────────────┴──────────┘
```

### Membuat Migrasi Baru

Untuk menambah migrasi baru, buat 2 file SQL di folder `db/migrations/`:

```
db/migrations/
├── 000001_create_tasks_table.up.sql      # Sudah ada
├── 000001_create_tasks_table.down.sql    # Sudah ada
├── 000002_add_priority_column.up.sql     # File baru (up)
└── 000002_add_priority_column.down.sql   # File baru (down)
```

Contoh isi file migrasi baru:

**`000002_add_priority_column.up.sql`**:

```sql
ALTER TABLE tasks ADD COLUMN priority VARCHAR(10) DEFAULT 'medium'
    CHECK (priority IN ('low', 'medium', 'high'));
```

**`000002_add_priority_column.down.sql`**:

```sql
ALTER TABLE tasks DROP COLUMN IF EXISTS priority;
```

Lalu jalankan: `go run main.go migrate`

## 📡 API Endpoints

### Health Check

```
GET /health
```

---

### 1. Create Task — `POST /tasks`

**Request Body:**

```json
{
  "title": "Belajar Go",
  "description": "Mempelajari bahasa pemrograman Go",
  "status": "pending",
  "due_date": "2026-04-01"
}
```

**Response (201):**

```json
{
  "message": "Task created successfully",
  "task": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Belajar Go",
    "description": "Mempelajari bahasa pemrograman Go",
    "status": "pending",
    "due_date": "2026-04-01",
    "created_at": "2026-03-17T08:00:00Z",
    "updated_at": "2026-03-17T08:00:00Z"
  }
}
```

---

### 2. Get All Tasks — `GET /tasks`

**Query Parameters:**

| Parameter | Tipe   | Deskripsi                             |
| --------- | ------ | ------------------------------------- |
| `status`  | string | Filter: `pending` atau `completed`    |
| `page`    | int    | Halaman (default: 1)                  |
| `limit`   | int    | Jumlah per halaman (default: 10)      |
| `search`  | string | Kata pencarian pada title/description |

**Response (200):**

```json
{
  "tasks": [...],
  "pagination": {
    "current_page": 1,
    "total_pages": 1,
    "total_tasks": 1
  }
}
```

---

### 3. Get Task by ID — `GET /tasks/:id`

**Response (200):** Objek task tunggal

---

### 4. Update Task — `PUT /tasks/:id`

**Request Body:** Sama seperti Create  
**Response (200):**

```json
{
  "message": "Task updated successfully",
  "task": {...}
}
```

---

### 5. Delete Task — `DELETE /tasks/:id`

**Response (200):**

```json
{
  "message": "Task deleted successfully"
}
```

---

### 6. Batch Create (Concurrent) — `POST /tasks/batch`

**Request Body:** Array dari task objects  
**Response (201):**

```json
{
  "message": "Batch creation completed",
  "created_count": 2,
  "error_count": 0,
  "tasks": [...]
}
```

## ❌ Error Responses

```json
{
  "error": "Validation failed",
  "details": {
    "title": "Title is required",
    "status": "Status must be 'pending' or 'completed'"
  }
}
```

| HTTP Code | Deskripsi                |
| --------- | ------------------------ |
| 400       | Bad Request / Validation |
| 404       | Task Not Found           |
| 500       | Internal Server Error    |

## ⚡ Concurrent Execution

1. **`GetAllTasks`** — Goroutines + channels untuk data fetch dan logging bersamaan, dengan `sync.WaitGroup`
2. **`CreateTasksConcurrently`** — Goroutines + `sync.Mutex` untuk batch create yang thread-safe
3. **Graceful Shutdown** — Server di goroutine terpisah, main goroutine menunggu sinyal interrupt

## 📝 Logging

Menggunakan [zerolog](https://github.com/rs/zerolog):

- Setiap HTTP request (method, path, status, duration)
- Setiap error di semua layer (repository, service, handler)
- Database events (koneksi, migrasi)
- Panic recovery

## 🧪 Testing dengan cURL

```bash
# Health check
curl http://localhost:8080/health

# Create task
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Belajar Go","description":"Mempelajari Go","status":"pending","due_date":"2026-04-01"}'

# Get all tasks
curl http://localhost:8080/tasks

# Get with filter & search
curl "http://localhost:8080/tasks?status=pending&search=belajar&page=1&limit=5"

# Update task (ganti {id} dengan UUID)
curl -X PUT http://localhost:8080/tasks/{id} \
  -H "Content-Type: application/json" \
  -d '{"title":"Updated","description":"Done","status":"completed","due_date":"2026-04-01"}'

# Delete task
curl -X DELETE http://localhost:8080/tasks/{id}

# Batch create
curl -X POST http://localhost:8080/tasks/batch \
  -H "Content-Type: application/json" \
  -d '[{"title":"A","description":"","status":"pending","due_date":"2026-04-01"},{"title":"B","description":"","status":"pending","due_date":"2026-04-02"}]'
```

## 📦 Dependencies

| Package                    | Deskripsi                           |
| -------------------------- | ----------------------------------- |
| `github.com/go-chi/chi/v5` | HTTP router yang ringan & idiomatic |
| `github.com/lib/pq`        | PostgreSQL driver untuk Go          |
| `github.com/rs/zerolog`    | Structured logging library          |

## 📄 Lisensi

MIT License
