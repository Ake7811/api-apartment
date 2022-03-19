package repository

import "apartment/api/models"

type BuildingRepository interface {
	FindAll() ([]models.Building, error)
	FindByID(uint64) (models.Building, error)
	Create(models.Building) (models.Building, error)
	Update(uint64, models.Building) (int64, error)
	Delete(uint64) (int64, error)

	CheckUniqueBuilding(models.Building) (bool, error)
}
