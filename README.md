# Notes App

## Table of Contents
- [Overview](#overview)
  - [Tools and Libraries](#tools-and-libraries)
  - [Architecture](#architecture)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Steps to Start the Server](#steps-to-start-the-server)
  - [Configure the Database (OPTIONAL)](#configure-the-database-optional)
  - [Migrations](#migrations)
  - [Seeding the Database](#seeding-the-database)
  - [Using the API with Postman](#using-the-api-with-postman)
- [Running Tests](#running-tests)
- [Teardown Instructions](#teardown-instructions)
- [API Endpoints](#api-endpoints)

This is a simple Notes application built with Go. The app allows users to create, retrieve, update, and delete notes. It follows a repository-style architecture for clean separation of concerns and includes comprehensive tests to ensure functionality.

There is currently no frontend for this app. I wanted to learn Go so implementing a UI is beyond the scope of this work. My goals for this app were to build a CRUD app that had some basic functionalities that are expected with backend projects, like testing, seeders, factories and a clean architecture. Implementing authorization and authentication is a future piece of this project.

---

## Overview

### Tools and Libraries
- **Go**: The primary programming language used for the application.
- **MySQL**: The main database for the app.
- **SQLite**: Database used for running API tests.
- **Gorilla Mux**: Handles HTTP requests.
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
- MySQL installed and running on your machine.
  - Ensure you have the MySQL credentials (e.g., username, password, and host).
  - The app will create a database called `mckinney_go_notes_db` during the migration process.

---

### Steps to Start the Server
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/go-notes.git
   cd go-notes
2. Initialize the database: `go run main.go migrate`
3. Start the server: `go run main.go`
4. The server will run on `http://localhost:8080`.

### Configure the Database (OPTIONAL)
- The code is setup to automatically install a MySQL DB with default env variables. If you would like to you can modify the default values in the config.go file.

### Migrations
- Run the following command to create the database and tables: `go run main.go --migrate`

### Seeding the Database
To populate the database with initial data, run the seeders: `go run main.go --seedDev`

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
