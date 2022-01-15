package repository

import "apartment/api/models"

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(uint32) (models.User, error)
	Create(models.User) (models.User, error)
	Update(uint32, models.User) (int64, error)
	Delete(uint32) (int64, error)
}
