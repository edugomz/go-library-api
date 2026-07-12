# Library API

## 🧠 Overview

A modular backend API for managing a digital library system, built with Go, Gin, and GORM.  
The project focuses on clean architecture, dependency separation, and production-style service structuring while remaining lightweight and educational.

---

## 🎯 Problem Statement

Building backend APIs at scale requires more than CRUD endpoints.  
This project explores:

- Layered architecture (handlers, services, repositories)
- Relational data modeling for real-world domains
- Dependency injection patterns in Go
- Containerized deployment workflows

The goal is to simulate how production Go services are structured without introducing unnecessary framework complexity.

---

## 🏗️ Architecture

* HTTP API built with Gin
* Layered architecture:
  * Handlers (HTTP layer)
  * Services (business logic)
  * Repositories (data access layer)
* PostgreSQL as the primary relational datastore
* GORM for ORM-based persistence
* Manual dependency wiring via application bootstrap layer
* Dockerized runtime environment

---

## ⚙️ Features

* User management (CRUD)
* Author and book catalog management
* Book reviews system
* Reading lists per user
* Relational integrity via foreign keys
* JWT authentication (register, login, Bearer token validation)
* Structured application bootstrap
* Environment-based configuration
* Container-ready deployment (Docker)

---

## 🛠️ Tech Stack

* Go (Gin framework)
* GORM (ORM)
* PostgreSQL
* Docker
* pgAdmin (development tooling)
* Clean layered architecture (service/repository pattern)

---

## 🔥 Design Highlights

* Explicit dependency injection without frameworks
* Separation of concerns across all layers
* Relational schema design with strong consistency guarantees
* Modular route registration (API versioning ready)
* Cloud-ready containerized deployment model

---

## 📊 Data Model

Core entities:

* Users
* Authors
* Books
* Reviews
* Reading Lists

Relationships are fully normalized and enforced at the database level using foreign keys.

---

## 🚀 Roadmap

* [x] Replace AutoMigrate with versioned migrations
* [x] Add authentication (JWT-based)
* [ ] Introduce caching layer (Redis)
* [x] Add integration + unit test suite
* [ ] Add API pagination and filtering
* [ ] Deploy to Cloud Run (GCP)
* [ ] Add observability (logging, metrics)

---

## 🧪 Local Setup

### Prerequisites

```bash
# Start infrastructure
docker-compose up -d

# Install dependencies
go mod tidy

# Run database migrations (one-off, run once per schema change)
go run ./cmd/api --migrate-only
# or: make migrate

#Run the application
go run ./cmd/api/main.go

```

### Migrations

Migrations no longer run automatically when the server boots.
On platforms like Cloud Run, a normal deploy or scale-up can start several instances at once, and concurrent migration runs against the same database can race.

Instead, run migrations as an explicit, separate step before starting (or updating) the server:

```bash
go run ./cmd/api --migrate-only
# or: make migrate
```

In production this should run as a single one-off job (e.g. a Cloud Run Job) before the new revision receives traffic, not as part of every instance's boot.
The same container image works for this: override its command with `./server --migrate-only` instead of building a separate image.

When deploying with `docker-compose.prod.yml`, run the profile-gated `migrate` service once before `up`:

```bash
docker-compose -f docker-compose.prod.yml --profile migrate run --rm migrate
docker-compose -f docker-compose.prod.yml up -d
```
