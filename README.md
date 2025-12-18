# go-api-auth-jwt — JWT Authentication (Gin + GORM + PostgreSQL)

โปรเจคนี้สอนทำระบบล็อกอินด้วย **JWT (Access Token + Refresh Token)** แบบใช้งานจริง พร้อม:
- Register / Login
- JWT Middleware (Protect Routes)
- Refresh Token
- Logout (Revoke Refresh Token)
- PostgreSQL + GORM
- Config ผ่าน .env (Viper)

> เหมาะสำหรับใช้เป็นฐานของระบบ API ขนาดกลาง–ใหญ่ และต่อยอด Microservices

---

## Tech Stack
- Go
- Gin
- GORM
- PostgreSQL
- JWT (Access / Refresh)
- bcrypt
- Viper
- Docker / Docker Compose

---

## API Endpoints

### Public
| Method | Path | Description |
|------|------|-------------|
| POST | /auth/register | สมัครสมาชิก |
| POST | /auth/login | ล็อกอิน |
| POST | /auth/refresh | ต่ออายุ access token |
| POST | /auth/logout | ออกจากระบบ |

### Protected
| Method | Path | Description |
|------|------|-------------|
| GET | /me | ดูข้อมูลผู้ใช้ |

---

## Project Structure

```bash
go-api-auth-jwt/
├─ cmd/api/main.go
├─ internal/
│  ├─ config
│  ├─ db
│  ├─ models
│  ├─ repositories
│  ├─ services
│  ├─ middlewares
│  ├─ handlers
│  └─ routes
├─ docker-compose.yml
├─ .env.example
└─ README.md
```

---

## Quick Start

### 1) Clone & Install
```bash
go mod tidy
```

### 2) Setup ENV
```bash
cp .env.example .env
```

### 3) Run PostgreSQL
```bash
docker compose up -d
```

### 4) Run API
```bash
go run ./cmd/api
```

Server: http://localhost:8080

---

## ENV Example

```env
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_auth

JWT_ACCESS_SECRET=change-me-access
JWT_REFRESH_SECRET=change-me-refresh

ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_DAYS=7
```

---

## Security Concept
- Hash password ด้วย bcrypt
- แยก Access / Refresh Secret
- Access Token อายุสั้น
- Refresh Token เก็บแบบ hash ใน DB
- Logout = revoke refresh token

---

## Next Improvements
- Refresh Token Rotation
- Role / Permission
- Email Verification
- Rate Limit
- Unit Test
- Observability (Log / Metrics)

---

## Author
Chinawat Daochai  
Course: Mastering Go API Development
# 011-go-api-auth-jwt
