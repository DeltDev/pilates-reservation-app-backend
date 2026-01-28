# Pilates Reservation System (Backend)
This is the backend service for the ZenPilates Reservation App. It handles court availability, and booking logic. Built with Go (Golang), Gin, and PostgreSQL.

This backend is live on Railway: https://pilates-reservation-app-backend-production.up.railway.app
## Tech Stack
Language: Go 1.25,6

Framework: Gin Web Framework

Database: PostgreSQL 18

Driver: pgx/v5

Migrations: golang-migrate

### 1. Prerequisites
Make sure you have these installed before running this backend into your local machine

Go (1.25.6)

Docker Desktop (to get the PostgreSQL 18 Image)

Golang-Migrate (4.19.1)

make (4.3)

### 2. Initialize PostgreSQL image

```make postgresinit```

or

```docker run --name postgres18 -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:18-alpine```

### 3. Create Database

```make createdb```

or

```docker exec -it postgres18 createdb --username=root --owner=root pilates```

### 4. Run main.go

```go run cmd/server/main.go```

### Other commands
See Makefile file in this repository to find more commands.