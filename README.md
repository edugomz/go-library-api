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

* [ ] Replace AutoMigrate with versioned migrations
* [ ] Add authentication (JWT-based)
* [ ] Introduce caching layer (Redis)
* [ ] Add integration + unit test suite
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

#Run the application
go run ./cmd/api/main.go

```
