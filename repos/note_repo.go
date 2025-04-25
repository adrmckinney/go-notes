package repos

import (
	"database/sql"
	"fmt"

	"github.com/adrmckinney/go-notes/models"
)

type NoteRepo struct {
	DB *sql.DB
}

func (r *NoteRepo) GetNoteById(id string) (models.Note, error) {
	var note models.Note
	err := r.DB.QueryRow("SELECT * FROM notes WHERE id = ?", id).Scan(
		&note.ID, &note.Title, &note.Content, &note.Added, &note.Modified,
	)
	if err != nil {
		return models.Note{}, err
	}
	return note, nil
}

func (r *NoteRepo) GetNotes() ([]models.Note, error) {
	rows, err := r.DB.Query("SELECT * FROM notes")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()
	var notes []models.Note
	for rows.Next() {
		var note models.Note
		if err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.Added, &note.Modified); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		notes = append(notes, note)
	}

	// Check for errors encountered during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during row iteration: %w", err)
	}

	return notes, nil
}

func (r *NoteRepo) CreateNote(note models.Note) (int64, error) {
	var result sql.Result
	var err error

	// Check if Added and Modified are empty
	if note.Added == "" || note.Modified == "" {
		// Use the database defaults for added and modified
		result, err = r.DB.Exec(
			"INSERT INTO notes (title, content) VALUES (?, ?)",
			note.Title, note.Content,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to insert note with defaults: %w", err)
		}
	} else {
		// Explicitly insert added and modified values
		result, err = r.DB.Exec(
			"INSERT INTO notes (title, content, added, modified) VALUES (?, ?, ?, ?)",
			note.Title, note.Content, note.Added, note.Modified,
		)
		if err != nil {
			return 0, fmt.Errorf("failed to insert note with explicit dates: %w", err)
		}
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve last inserted ID: %w", err)
	}

	return id, nil
}

func (r *NoteRepo) UpdateNote(id string, note models.Note) error {
	query := "UPDATE notes SET "
	args := []interface{}{} // The first {} defines the slice type, and the second {} initializes it with no elements.
	if note.Title != "" {
		query += "title = ?, "
		args = append(args, note.Title)
	}

	if note.Content != "" {
		query += "content = ? "
		args = append(args, note.Content)
	}

	query = query[:len(query)-2] // remove the trailing comma and space
	query += " WHERE id = ?"
	args = append(args, id)

	_, err := r.DB.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update note %v", err)
	}
	return nil
}

func (r *NoteRepo) DeleteNote(id string) (map[string]string, error) {
	res, err := r.DB.Exec("DELETE FROM notes where id = ?", id)
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to delete note: %v", err)
	}

	// Check if any rows were affected by the delete
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return map[string]string{}, fmt.Errorf("failed to retrieve delete result: %v", err)
	}

	if rowsAffected == 0 {
		return map[string]string{}, fmt.Errorf("note not found: %v", err)
	}
	return map[string]string{
		"message": "Note successfully deleted",
	}, nil
}
