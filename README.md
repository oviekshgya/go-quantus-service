# 🚀 Go Quantus Service

> A scalable Go microservice using Gin, PostgreSQL, Redis, RabbitMQ, and Docker.

## 📚 API Documentation

📄 [Open API Docs (via Apidog)](https://dc41pw3e0w.apidog.io)

---

## ⚙️ Tech Stack

- 🧠 **Go** (v1.24)
- 🔥 **Gin Gonic** Framework
- 🐘 **PostgreSQL** (15-alpine)
- 🚀 **Redis** (7-alpine)
- 📩 **RabbitMQ** (3-management)
- 🐳 **Docker** & Docker Compose
- 🧾 **Custom Logger Worker Middleware**

---

## 🛠️ Setup & Run

### 1. 🔧 Prerequisites

- Docker & Docker Compose
- `.env` file in root directory with:

```env
SERVICE_PORT=7003
DB_HOST=postgres
DB_PORT=5432
DB_USER=youruser
DB_PASSWORD=yourpass
DB_NAME=yourdb
REDIS_PORT=6379
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASS=guest
```

---

# Build dan jalankan semua service
````cmd
docker-compose up --build -d
````

---

# 🧪 UserController Unit Tests

Dokumen ini menjelaskan pengujian unit untuk `UserController` pada proyek `go-quantus-service`. Pengujian dilakukan dengan menggunakan:

- [Gin](https://github.com/gin-gonic/gin) sebagai web framework
- [Testify](https://github.com/stretchr/testify) untuk assertion
- [GoMock](https://github.com/golang/mock) untuk mocking service
- [net/http/httptest](https://pkg.go.dev/net/http/httptest) untuk simulasi HTTP server


```go
cd test
go test -v .
```

---
# 🧪 Async Log Bonus:

```go
func StartLogWorker(db *gorm.DB, redisClient *redis.RedisClient, batchSize int, interval time.Duration) {
	ticker := time.NewTicker(interval)

	go func() {
		for range ticker.C {
			queueLen, err := redisClient.C.LLen("log_queue").Result()
			if err != nil {
				fmt.Println("Redis LLEN error:", err)
				continue
			}
			if queueLen >= int64(batchSize) {
				logStrings, err := redisClient.C.LRange("log_queue", 0, int64(batchSize-1)).Result()
				if err != nil {
					fmt.Println("Redis LRANGE error:", err)
					continue
				}

				var logs []entities.LogEntry
				for _, item := range logStrings {
					var log entities.LogEntry
					if err := json.Unmarshal([]byte(item), &log); err == nil {
						logs = append(logs, log)
					}
				}

				if len(logs) > 0 {
					if err := db.Create(&logs).Error; err != nil {
						fmt.Println("DB insert error:", err)
						continue 
					}

					if err := redisClient.C.LTrim("log_queue", int64(batchSize), -1).Err(); err != nil {
						fmt.Println("Redis LTRIM error:", err)
					}

				}
			}
		}
	}()
}
```

Middleware:
````go
func RequestLogger(logController controller.LogControllerinterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Read body
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		// Proses request
		c.Next()

		// After request
		status := c.Writer.Status()
		headers, _ := json.Marshal(c.Request.Header)

		log := entities.LogEntry{
			IPAddress: c.ClientIP(),
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Headers:   string(headers),
			Body:      string(bodyBytes),
			Response:  strconv.Itoa(status),
			Status:    status,
			CreatedAt: start,
		}

		// Ambil Redis dari controller
		_ = logController.GetDependencies().Redis.PushLogToQueue("log_queue", log)
	}
}

````
---

## 📌 Endpoint API

Semua endpoint base URL:


---

### 🔐 Auth & Health Check

| Method | Endpoint   | Auth       | Deskripsi             |
|--------|------------|------------|------------------------|
| GET    | `/ping`    | -          | Cek status service     |

---

### 👤 User Endpoints (`/users`)

| Method | Endpoint              | Auth        | Deskripsi                    |
|--------|-----------------------|-------------|------------------------------|
| POST   | `/users`              | Basic Auth  | Register user baru           |
| POST   | `/users/login`        | Basic Auth  | Login user, return JWT       |
| GET    | `/users`              | JWT         | Get user saat ini (by token) |
| GET    | `/users/:user_id`     | Basic Auth  | Get user by ID               |
| PUT    | `/users/:user_id`     | JWT         | Update user by ID            |
| DELETE | `/users/:user_id`     | JWT         | Hapus user by ID             |

---

### 📄 Content Endpoints (`/content`)

| Method | Endpoint                    | Auth | Deskripsi                         |
|--------|-----------------------------|------|-----------------------------------|
| POST   | `/content`                  | JWT  | Tambah konten baru                |
| GET    | `/content`                  | JWT  | Ambil semua konten                |
| GET    | `/content/:content_id`      | JWT  | Ambil konten berdasarkan ID       |
| PUT    | `/content/:content_id`      | JWT  | Update konten                     |
| DELETE | `/content/:content_id`      | JWT  | Soft delete konten                |
| DELETE | `/content/clean/:content_id`| JWT  | Hard delete konten (bersih total) |
| PATCH  | `/content/clean/:content_id`| JWT  | Update konten setelah bersih      |
