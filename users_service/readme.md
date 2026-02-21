# 🚀 User Service – Go Microservice Example

A minimal production-style REST API written in **Go**, demonstrating:

✅ Go modules  
✅ Clean project structure  
✅ Multiple packages  
✅ Middleware  
✅ Structured logging  
✅ Environment-based configuration  
✅ Docker multi-stage build  
✅ Small container image  

---

# 📦 Project Overview

**Service Name:** `user-service`  
**Tech Stack:** Go (net/http), Docker  
**Purpose:** Learn API development & microservice structure in Go  

---

# 🧰 Prerequisites

Before running the project, install:

- **Go** (1.20+ recommended)  
  https://go.dev/dl/

- **Docker** (optional, for containerization)  
  https://www.docker.com/products/docker-desktop/

Verify installations:

```bash
go version
docker --version
````

---

# 🚀 Setup & Installation

## 1️⃣ Clone or Create Project

```bash
mkdir user-service
cd user-service
```

## 2️⃣ Initialize Go Module

```bash
go mod init user-service
```

This creates:

```
go.mod → Dependency management file
```

---

# 📁 Project Structure

```
user-service/
│── go.mod
│── main.go
│
├── config/
│   └── config.go
│
├── handlers/
│   └── user_handler.go
│
├── models/
│   └── user.go
│
├── middleware/
│   └── logging.go
│
├── routes/
│   └── routes.go
│
├── utils/
│   └── logger.go
│
└── Dockerfile
```

---

# 🧠 File-by-File Explanation

---

## 🔹 `main.go` – Entry Point

**Responsibilities:**

* Load configuration
* Initialize logger
* Setup routes
* Start HTTP server

```go
func main() {
    cfg := config.LoadConfig()
    logger := utils.NewLogger()

    router := routes.SetupRoutes(logger)

    addr := fmt.Sprintf(":%s", cfg.Port)
    logger.Info("Starting server on " + addr)

    http.ListenAndServe(addr, router)
}
```

---

## 🔹 `config/config.go` – Configuration Management

**Responsibilities:**

* Read environment variables
* Provide defaults

```go
type Config struct {
    Port string
}
```

Loads:

```
PORT → Server port (default: 8080)
```

---

## 🔹 `models/user.go` – Data Structures

Defines API response objects.

```go
type User struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}
```

---

## 🔹 `handlers/user_handler.go` – Business Logic

Contains HTTP handlers:

* `/health`
* `/users`

Example:

```go
func GetUsers(w http.ResponseWriter, r *http.Request)
```

Returns JSON list of users.

---

## 🔹 `routes/routes.go` – Route Definitions

Centralizes endpoint registration.

```go
mux.HandleFunc("/health", handlers.HealthHandler)
mux.HandleFunc("/users", handlers.GetUsers)
```

Wraps middleware.

---

## 🔹 `middleware/logging.go` – Middleware Layer

Logs request details:

✔ Method
✔ Path
✔ Execution time

```go
logger.Info(r.Method + " " + r.URL.Path)
```

---

## 🔹 `utils/logger.go` – Logging Utility

Simple structured logger:

```
[INFO]
[ERROR]
```

---

## 🔹 `Dockerfile` – Containerization

Uses **multi-stage build** for small image.

Stages:

1️⃣ Build binary
2️⃣ Run on Alpine Linux

---

# ▶️ Running the Application Locally

```bash
go run main.go
```

Server starts at:

```
http://localhost:8080
```

---

# 🔍 Available Endpoints

## ✅ Health Check

```http
GET /health
```

Response:

```json
{"status":"ok"}
```

---

## 👥 Get Users

```http
GET /users
```

Response:

```json
[
  {"id":"1","name":"Jay","age":25},
  {"id":"2","name":"Rahul","age":30}
]
```

---

# ⚙️ Environment Variables

| Variable | Description | Default |
| -------- | ----------- | ------- |
| PORT     | Server port | 8080    |

Example:

```bash
PORT=9090 go run main.go
```

---

# 🐳 Docker Support

---

## 🏗 Build Docker Image

```bash
docker build -t user-service .
```

---

## ▶️ Run Container

```bash
docker run -p 8080:8080 user-service
```

Access API:

```
http://localhost:8080/health
```

---

# 📦 Why Multi-Stage Build?

Without multi-stage:

❌ Image size ~800MB+

With multi-stage + Alpine:

✅ Image size ~10–20MB 😎

Benefits:

✔ Faster deployment
✔ Less storage
✔ Better for microservices

---

# 🧪 Testing with Curl

```bash
curl http://localhost:8080/health
curl http://localhost:8080/users
```

---

# 🎯 Learning Objectives

By building this project you learn:

✅ Go modules (`go mod init`)
✅ Package system
✅ net/http server
✅ Handlers & routing
✅ Middleware
✅ Logging
✅ Environment config
✅ Dockerizing Go apps
✅ Multi-stage builds

---

# 🚀 Next Improvements (Optional)

You can extend this service with:

🔥 Gorilla Mux / Chi router
🔥 PostgreSQL / MongoDB
🔥 JWT Authentication
🔥 Graceful shutdown
🔥 Prometheus metrics
🔥 Request validation
🔥 Dependency injection


```

---

next we can build:

✅ Go + PostgreSQL microservice  
✅ Auth service (JWT)  
✅ API Gateway  
✅ Multi-container Docker Compose setup  
```
