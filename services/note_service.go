package services

import (
	"net/http"

	"github.com/adrmckinney/go-notes/errs"
	"github.com/adrmckinney/go-notes/models"
	"github.com/adrmckinney/go-notes/repos"
	"gorm.io/gorm"
)

type NoteService struct {
	NoteRepo repos.NoteRepo
}

func (s *NoteService) GetNoteById(authUserId uint, noteId uint) (models.Note, error) {
	note, err := s.NoteRepo.GetNoteById(noteId)

	if err != nil {
		return models.Note{}, errs.NewHTTPError(http.StatusNotFound, "Note not found")
	}
	if authUserId != note.UserID {
		return models.Note{}, errs.NewHTTPError(http.StatusUnauthorized, "Note does not belong to current user")
	}

	return note, nil
}

func (s NoteService) GetNotes(authUserId uint) ([]models.Note, error) {

	notes, err := s.NoteRepo.GetNotes(authUserId)
	if err != nil {
		return []models.Note{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to fetch notes")
	}
	return notes, nil
}

func (s NoteService) CreateNote(authUserId uint, payload models.Note) (models.Note, error) {
	if payload.Title == "" || payload.Content == "" {
		return models.Note{}, errs.NewHTTPError(http.StatusBadRequest, "Title and Content are required")
	}
	payload.UserID = authUserId

	createdNote, err := s.NoteRepo.CreateNote(payload)
	if err != nil {
		return models.Note{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to create note")
	}

	return createdNote, nil
}

func (s *NoteService) UpdateNote(authUserId uint, noteId uint, payload map[string]any) (models.Note, error) {
	// Validate note belongs to user
	_, err := s.GetNoteById(authUserId, noteId)
	if err != nil {
		return models.Note{}, err
	}
	updatedNote, err := s.NoteRepo.UpdateNote(noteId, payload)

	if err != nil {
		// Check for GORM's record not found error
		if err == gorm.ErrRecordNotFound ||
			(err.Error() != "" && (err.Error() == "note not found" ||
				err.Error() == "note not found: record not found")) {
			return models.Note{}, errs.NewHTTPError(http.StatusNotFound, "Note not found")
		}
		return models.Note{}, errs.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return updatedNote, nil
}

func (s *NoteService) DeleteNote(authUserId uint, noteId uint) (map[string]string, error) {
	// Validate note belongs to user
	_, err := s.GetNoteById(authUserId, noteId)
	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	err = s.NoteRepo.DeleteNote(noteId)

	if err != nil {
		return map[string]string{}, errs.NewHTTPError(http.StatusInternalServerError, "Failed to delete note")
	}
	return map[string]string{
		"message": "Note successfully deleted",
	}, nil
}
