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
  - [Running Tests](#running-tests)
  - [Teardown Instructions](#teardown-instructions)
  - [API Endpoints](#api-endpoints)
    - [Notes](#notes)

This is a simple Notes application built with Go. The app allows users to create, retrieve, update, and delete notes. It follows a repository-style architecture for clean separation of concerns and includes comprehensive tests to ensure functionality.

There is currently no frontend for this app. I wanted to learn Go so implementing a UI is beyond the scope of this work. My goals for this app were to build a CRUD app that had some basic functionalities that are expected with backend projects, like testing, seeders, factories and a clean architecture. Implementing authorization and authentication is a future piece of this project.

---

## Overview

### Tools and Libraries
- **Go**: The primary programming language used for the application.
- **GORM**: The ORM (Object Relational Mapper) used for database interactions, migrations, and model management.
- **MySQL**: The main database for the app.
- **SQLite**: Database used for running API tests.
- **Gorilla Mux**: Handles HTTP requests.
- **golang-migrate**: For running database migrations in a controlled, versioned manner.
- **Testing**:
  - Built-in Go testing framework (`testing` package).
  - `httptest` for simulating HTTP requests and responses.
- **Factories**: Custom factories for generating test and seeder data.
- **Seeders**: Functions to populate the database with initial data for testing and development.

### Architecture
The app follows a **repository-style architecture**:
- **Handlers**: Handle HTTP requests and responses. They interact with repositories to perform database operations.
- **Repositories**: Encapsulate database logic and provide methods for CRUD operations.
- **Models**: Define the structure of the data (e.g., `Note`).
- **Test Database Initialization**: Centralized logic for setting up the SQLite database.

### Tests
The app includes comprehensive tests for all major functionalities:
- **CreateNote**: Tests the creation of a new note.
- **UpdateNote**: Tests updating an existing note.
- **GetNote**: Tests retrieving a specific note by ID.
- **GetNotes**: Tests retrieving all notes.
- **DeleteNote**: Tests deleting a note by ID.
- **Test Utilities**:
  - `CleanUpDatabases`: Ensures the database is cleaned up after each test.
  - `Factories`: Generate test data for notes.

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
  make migrate-up
  ```
- **To rollback the last migration:**
  ```bash
  make migrate-down
  ```
- **To rollback multiple steps (e.g., 2 steps):**
  ```bash
  make migrate-steps steps=2
  ```
- **To check the current migration version:**
  ```bash
  make migrate-version
  ```
- You can override the database host if needed (e.g., inside Docker):
  ```bash
  DB_HOST=db make migrate-up
  ```
- See the `Makefile` for more migration commands and options.

### Seeding the Database
- To populate the database with initial data, run the seeders: `go run main.go --seedDev`

### Using the API with Postman
Once the server is running, you can use Postman to interact with the API. Below are the available endpoints:

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
### Notes
- Get All Notes: `GET /notes`
- Get Note by ID: `GET /notes/{id}`
- Create Note: `POST /notes`
- Update Note: `PUT /notes`
- Delete Note: `DELETE /notes/{id}`
