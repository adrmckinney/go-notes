# Notes App

## Table of Contents
- [Notes App](#notes-app)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
    - [Tools and Libraries](#tools-and-libraries)
    - [Architecture](#architecture)
    - [Tests](#tests)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Get the app](#get-the-app)
    - [Using Docker to Run the Project](#using-docker-to-run-the-project)
    - [Migrations](#migrations)
    - [Seeding the Database](#seeding-the-database)
    - [Using the API with Postman](#using-the-api-with-postman)
      - [Users](#users)
      - [Notes](#notes)
  - [Running Tests](#running-tests)
  - [Teardown Instructions](#teardown-instructions)
  - [API Endpoints](#api-endpoints)
    - [Users](#users)
    - [Notes](#notes-1)

This is a simple Notes application built with Go. It supports user authentication and authorization, allowing each user to manage their own notes. The app follows a clean, layered architecture with separation of concerns and includes comprehensive tests. Features include signup/signin, CRUD operations for notes, and database seeding with a demo user for testing.

There is currently no frontend for this app. I wanted to learn Go so implementing a UI is beyond the scope of this work. My goals for this app were to build a CRUD app that had some basic functionalities that are expected with backend projects, like AuthMiddleware for authentication, authorization, testing, seeders, factories and a clean architecture.

---

## Overview

### Tools and Libraries
- **Go**: The primary programming language used for the application.
- **GORM**: The ORM (Object Relational Mapper) used for database interactions, migrations, and model management.
- **MySQL**: The main database for the app.
- **SQLite**: Database used for running API tests.
- **Gorilla Mux**: Handles HTTP requests.
- **golang-migrate**: For running database migrations in a controlled, versioned manner.
- **JWT** (github.com/golang-jwt/jwt): Used for authentication tokens.
- **Validator** (go-playground/validator): Struct validation for request payloads.
- **Godotenv**: Loads environment variables from .env files.
- **Testing**:
  - Built-in Go testing framework (`testing` package).
  - `httptest` for simulating HTTP requests and responses.
- **Factories**: Custom factories for generating test and seeder data.
- **Seeders**: Functions to populate the database with initial data for testing and development.

### Architecture
The application uses a layered architecture:
-	**HTTP Layer (http/):**
    - Contains route registration and HTTP handlers for both notes and users.
    - Includes middleware for authentication using JWTs (AuthMiddleware).
- **Service Layer (services/):**
    - Contains business logic for users and notes.
    - Handles authorization, input validation, and calling repository methods.
- **Repository Layer (repositories/):**
    - Handles all interactions with the database.
    - Implements Create, Read, Update, Delete (CRUD) for users and notes.
- **Models (models/):**
    - Defines the structure of entities (e.g., User, Note).
- **Factories & Seeders (seeder/):**
    - Generate test and development data.
    - The first user seeded is a demo user:
    ```json
      {
        "username": "demoUser",
        "password": "password"
      }
    ```
- **Authentication & Authorization:**
    - JWT-based authentication.
    - Middleware ensures users can only access their own notes.
- **Test Database Initialization**: Centralized logic for setting up the SQLite database.

### Tests
The app includes API tests covering:
- CRUD operations for notes.
- User signup and signin.
- Authorization logic.
- Test utilities include:
    - CleanUpTables for DB cleanup between tests.
    - Factories for generating test users and notes.

---

## Getting Started

### Prerequisites
- Go installed on your machine (version 1.18 or higher).
- Docker installed (if using Docker).
- Docker Compose installed (if using Docker).
- Postman or any other API testing tool.
---

### Get the app
- Clone it: `git clone https://github.com/adrmckinney/go-notes.git`

### Using Docker to Run the Project

1. **Install Docker**:
   - Ensure Docker and Docker Compose are installed on your machine.

2. **Build and Start the Containers**:
   - Run the following command to build and start the containers:
     ```bash
     docker compose up --build
     ```

3. **Access the Application**:
   - The application will be available at `http://localhost:8080`.

4. **Stop the Containers**:
   - To stop the containers, press `Ctrl+C` or run:
     ```bash
     docker-compose down
     ```

5. **Database Persistence**:
   - The MySQL database data is stored in a Docker volume (`db_data`) to ensure persistence across container restarts.

### Migrations
- This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database migrations.
- **To run migrations (apply all up migrations):**
  ```bash
  docker compose exec app make migrate-up
  ```
- **To rollback the last migration:**
  ```bash
  docker compose exec app make migrate-down
  ```
- **To rollback multiple steps (e.g., 2 steps):**
  ```bash
  docker compose exec app make migrate-steps steps=2
  ```
- **To check the current migration version:**
  ```bash
  docker compose exec app make migrate-version
  ```
- See the `Makefile` for more migration commands and options.

### Seeding the Database
- To populate the database with initial data, run the seeders: `docker compose exec app make seed-dev`
- This will create initial data.
- Demo user credentials:
  ```json
      {
        "username": "demoUser",
        "password": "password"
      }
    ```

### Using the API with Postman
Once the server is running, you can use Postman to interact with the API. Below are the available endpoints:

**All user and note routes (except /signup and /signin) require a valid Authorization: Bearer <token> header.**

#### Users
- Signup
    - `POST /signup`
    - Payload:
  ```json
  {
    "firstName": "Tom",
    "lastName": "Holland",
    "username": "demoUser",
    "password": "password",
    "confirmPassword": "password"
  }
  ```

- Signin
    - `POST /signin`
    - Payload:
  ```json
  {
    "username": "demoUser",
    "password": "password"
  }
  ```

- Update User
    - `PUT /signin`
    - Payload:
  ```json
  {
    "firstName": "UpdatedFirst",
    "lastName": "UpdatedLast",
    "password": "newPassword",
    "confirmPassword": "newPassword"
  }
  ```

#### Notes
- Get All Notes:
   - URL: http://localhost:8080/notes
   - Method: GET
   - Get Note by ID:

- Get Note
  - URL: http://localhost:8080/notes/{id}
  - Method: GET
  - Example: http://localhost:8080/notes/1

- Create Note:
  - URL: http://localhost:8080/notes
  - Method: POST
  - Payload (JSON):
```json
{
  "title": "Sample Note",
  "content": "This is a sample note."
}
```
- Update Note:
  - URL: http://localhost:8080/notes
  - Method: PUT
  - Payload (JSON):
```json
{
  "title": "Updated Title",
}
```
- Delete Note:
  - URL: http://localhost:8080/notes/{id}
  - Method: DELETE
  - Example: http://localhost:8080/notes/1


## Running Tests
To run the tests, use the following command: `go test -v ./tests/...`

## Teardown Instructions
To completely remove the database, run: `go run main.go --remove`

## API Endpoints
### Users
- Sign Up: `POST /signup`
- Sign In: `POST /signin`
- Update User: `PUT /user`

### Notes
- Get All Notes: `GET /notes`
- Get Note by ID: `GET /notes/{id}`
- Create Note: `POST /notes`
- Update Note: `PUT /notes`
- Delete Note: `DELETE /notes/{id}`
