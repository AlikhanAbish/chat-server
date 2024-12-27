# chat-serveк
# Project Overview
This project sets up a simple Go server that handles JSON-based HTTP requests (POST and GET), connects to a PostgreSQL database, and implements basic CRUD operations. It also includes optional HTML interface integration for testing.

# Dependencies
Go (version 1.18+)
PostgreSQL (for database connection)

# Go Packages:
github.com/lib/pq (PostgreSQL driver for Go)
gorm.io/gorm (ORM for Go)
github.com/golang-migrate/migrate (for database migrations)

# Install Dependencies: In your project directory, run:
go get github.com/lib/pq
go get gorm.io/gorm
go get -u github.com/golang-migrate/migrate/cmd/migrate

# Run the Server: In your project folder, run:
 Run the Server: In your project folder, run:
