package repos

import (
	"github.com/adrmckinney/go-notes/models"
	"gorm.io/gorm"
)

type UserTokenRepo struct {
	DB *gorm.DB
}

func (r *UserTokenRepo) StoreUserToken(userToken models.UserToken) error {
	result := r.DB.Create(&userToken)
	return result.Error
}

func (r *UserTokenRepo) DeleteUserToken(token string) error {
	return r.DB.Where("token = ?", token).Delete(&models.UserToken{}).Error
}

func (r *UserTokenRepo) TokenExists(token string) bool {
	err := r.DB.Where("token = ?", token).First(&models.UserToken{}).Error
	return err == nil
}
