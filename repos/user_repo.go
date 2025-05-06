package repos

import (
	"fmt"

	"github.com/adrmckinney/go-notes/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func (r *UserRepo) GetUserById(id uint) (models.User, error) {
	var user models.User
	result := r.DB.First(&user, id)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepo) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	result := r.DB.Where("username = ?", username).First(&user)
	return user, result.Error
}

func (r *UserRepo) GetUsers() ([]models.User, error) {
	var users []models.User
	result := r.DB.Find(&users)
	if result.Error != nil {
		return []models.User{}, result.Error
	}

	return users, nil
}

func (r *UserRepo) CreateUser(user models.User) (models.User, error) {
	result := r.DB.Create(&user)
	if result.Error != nil {
		return models.User{}, result.Error
	}
	return user, nil
}

func (r *UserRepo) UpdateUser(id uint, updated map[string]interface{}) (models.User, error) {
	var user models.User
	if result := r.DB.First(&user, id); result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return models.User{}, gorm.ErrRecordNotFound
		}
		return models.User{}, fmt.Errorf("user not found: %w", result.Error)
	}

	if err := r.DB.Model(&user).Updates(updated).Error; err != nil {
		return models.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	if err := r.DB.First(&user, id).Error; err != nil {
		fmt.Println("ERROR in FINAL", err)
		return models.User{}, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	return user, nil
}

func (r *UserRepo) DeleteUser(id uint) error {
	result := r.DB.Delete(&models.User{}, id)
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return result.Error
}
