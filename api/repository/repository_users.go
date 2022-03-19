package repository

import "apartment/api/models"

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(uint64) (models.User, error)
	Create(models.User) (models.User, error)
	Update(uint64, models.User) (int64, error)
	Delete(uint64) (int64, error)

	CheckUniqueUser(models.User) (bool, error)
}
