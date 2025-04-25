# Notes App

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
- **Database Initialization**: Centralized logic for setting up the SQLite database.

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
- SQLite installed (optional, as the app can use an in-memory database for testing).

### Steps to Start the Server
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/go-notes.git
   cd go-notes
2. Initialize the database: `go run main.go migrate`
3. Start the server: `go run main.go`
4. The server will run on `http://localhost:8080`.

### Seeding the Database
To populate the database with initial data, run the seeders: `go run main.go --seedDev`

### Running Tests
To run the tests, use the following command: `go test -v ./tests/...`

## API Endpoints
### Notes
- Get All Notes: `GET /notes`
- Get Note by ID: `GET /notes/{id}`
- Create Note: `POST /notes`
- Update Note: `PUT /notes`
- Delete Note: `DELETE /notes/{id}`

## Repository-Style Architecture
### Handlers
- Handle HTTP requests and responses.
- Example: `CreateNote`, `GetNote`, `DeleteNote`.
### Repositories
- Encapsulate database logic.
- Example: `NoteRepo` provides methods like `CreateNote`, `GetNoteById`, `DeleteNote`.
### Models
- Define the structure of the data.
- Example: `Note` model includes fields like `ID`, `Title`, `Content`, `Added`, and `Modified`.
