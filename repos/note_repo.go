package repos

import (
	"fmt"

	"github.com/adrmckinney/go-notes/models"
	"gorm.io/gorm"
)

type NoteRepo struct {
	DB *gorm.DB
}

func (r *NoteRepo) GetNoteById(id uint) (models.Note, error) {
	var note models.Note
	result := r.DB.First(&note, id)
	if result.Error != nil {
		return models.Note{}, result.Error
	}
	return note, nil
}

func (r *NoteRepo) GetNotes(userId uint) ([]models.Note, error) {
	var notes []models.Note
	result := r.DB.Where("user_id = ?", userId).Find(&notes)
	if result.Error != nil {
		return []models.Note{}, result.Error
	}

	return notes, nil
}

func (r *NoteRepo) CreateNote(note models.Note) (models.Note, error) {
	result := r.DB.Create(&note)
	if result.Error != nil {
		return models.Note{}, result.Error
	}
	return note, nil
}

func (r *NoteRepo) UpdateNote(id uint, updated map[string]interface{}) (models.Note, error) {
	var note models.Note
	if result := r.DB.First(&note, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.Note{}, gorm.ErrRecordNotFound
		}
		return models.Note{}, fmt.Errorf("note not found: %w", result.Error)
	}

	if err := r.DB.Model(&note).Updates(updated).Error; err != nil {
		return models.Note{}, fmt.Errorf("failed to update note: %w", err)
	}

	if err := r.DB.First(&note, id).Error; err != nil {
		fmt.Println("ERROR in FINAL", err)
		return models.Note{}, fmt.Errorf("failed to fetch updated note: %w", err)
	}

	return note, nil
}

func (r *NoteRepo) DeleteNote(id uint) error {
	result := r.DB.Delete(&models.Note{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("note not found")
	}
	return result.Error
}
