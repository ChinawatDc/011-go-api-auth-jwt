# go-api-auth-jwt — JWT Authentication (Gin + GORM + PostgreSQL)

โปรเจคนี้สอนทำระบบล็อกอินด้วย **JWT (Access Token + Refresh Token)** แบบใช้งานจริง พร้อม:

- Register / Login
- JWT Middleware (Protect Routes)
- Refresh Token (ออก Access ใหม่)
- Logout (Revoke Refresh Token)
- เก็บผู้ใช้ + refresh token ใน PostgreSQL (ผ่าน GORM)
- Config ผ่าน `.env` (ใช้ Viper)

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

## API Overview

### Public

- POST /auth/register
- POST /auth/login
- POST /auth/refresh
- POST /auth/logout

### Protected

- GET /me

---

## Project Structure

```
go-api-auth-jwt/
├─ cmd/api/main.go
├─ internal/
│  ├─ config/
│  ├─ db/
│  ├─ models/
│  ├─ repositories/
│  ├─ services/
│  ├─ middlewares/
│  ├─ handlers/
│  └─ routes/
├─ .env.example
├─ docker-compose.yml
├─ go.mod
└─ README.md
```

---

## Environment Config (.env.example)

```
APP_NAME=go-api-auth-jwt
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_auth

JWT_ISSUER=go-api-auth-jwt
JWT_ACCESS_SECRET=super-access-secret-change-me
JWT_REFRESH_SECRET=super-refresh-secret-change-me

ACCESS_TOKEN_MINUTES=15
REFRESH_TOKEN_DAYS=7
```

---

## Docker Compose

```
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: go_auth
    ports:
      - "5432:5432"
```

---

## Security Design

- Password hash ด้วย bcrypt
- Access token อายุสั้น
- Refresh token อายุยาว เก็บแบบ hash ใน DB
- แยก secret access / refresh
- ตรวจ revoke / expire ทุกครั้ง

---

## Example cURL

### Register

```
curl -X POST http://localhost:8080/auth/register -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"123456"}'
```

### Login

```
curl -X POST http://localhost:8080/auth/login -H "Content-Type: application/json" -d '{"email":"test@example.com","password":"123456"}'
```

### Protected Route

```
curl http://localhost:8080/me -H "Authorization: Bearer ACCESS_TOKEN"
```

---

## Next Step

- Refresh Token Rotation
- Rate Limit
- Email Verification
- RBAC
- Logging / Monitoring
- Unit Test

---

Author: Chinawat Daochai
