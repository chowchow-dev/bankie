# Bankie

Simple bank service written in Golang

## Checklist

- [x] Database tasks
  - [x] Database + API design
  - [x] Database migration
  - [x] Auto generate CRUD Golang code from SQL
  - [x] Implement database transaction
  - [x] Unit test for database CRUD
  - [x] Infra + Github Actions + Config
- [ ] RESTful API
  - [x] Implement APIs
  - [ ] Mock database for testing http API

## Setup local development

### Tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

  ```bash
  brew install golang-migrate
  ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

  ```bash
  brew install sqlc
  ```

### Infra

- Start postgres container:

  ```bash
  make postgres
  ```

- Create simple_bank database:

  ```bash
  make createdb
  ```

- Run db migration up all versions:

  ```bash
  make migrateup
  ```

- Run db migration down all versions:

  ```bash
  make migratedown
  ```

### How to run

- Add `local.env` file:

  ```env
    DB_DRIVER=
    DB_SOURCE=
    SERVER_ADDRESS=
  ```

- Run server:

  ```bash
  make server
  ```

- Run test:

  ```bash
  make test
  ```
