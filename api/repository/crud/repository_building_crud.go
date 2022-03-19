package crud

import (
	"apartment/api/models"
	"apartment/api/utils/channels"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type repositoryBuildingCRUD struct {
	db *gorm.DB
}

func NewRepositoryBuildingCRUD(db *gorm.DB) *repositoryBuildingCRUD {
	return &repositoryBuildingCRUD{db}
}
func (r *repositoryBuildingCRUD) FindAll() ([]models.Building, error) {
	var err error
	buildings := []models.Building{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		//Get only active user
		err = r.db.Debug().Model(&models.Building{}).Where("building_status > 0 ").Limit(100).Find(&buildings).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return buildings, nil
	}
	return nil, err
}
func (r *repositoryBuildingCRUD) FindByID(uid uint64) (models.Building, error) {
	var err error
	buildings := models.Building{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.User{}).Where("id = ?", uid).Take(&buildings).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return buildings, nil
	}
	if gorm.IsRecordNotFoundError(err) {
		return models.Building{}, errors.New("building not found")
	}
	return models.Building{}, err
}
func (r *repositoryBuildingCRUD) Create(sourceBuilding models.Building) (models.Building, error) {
	var err error
	var found bool
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		//Check Unique
		found, err = r.CheckUniqueBuilding(sourceBuilding)
		if found { // found same data in db
			ch <- false
			return
		}
		err = r.db.Debug().Model(&models.Building{}).Create(&sourceBuilding).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return sourceBuilding, nil
	}
	return models.Building{}, err
}
func (r *repositoryBuildingCRUD) Update(uid uint64, sourceBuilding models.Building) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.Building{}).Where("id = ?", uid).Take(&models.Building{}).UpdateColumns(
			map[string]interface{}{
				"building_name":   sourceBuilding.Building_name,
				"remark":          sourceBuilding.Remark,
				"building_status": sourceBuilding.Building_status,
				"updated_at":      time.Now(),
			},
		)
		ch <- true
	}(done)

	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}
		return rs.RowsAffected, nil
	}
	return 0, rs.Error
}
func (r *repositoryBuildingCRUD) Delete(uid uint64) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		//Delete is update building_status = 0
		rs = r.db.Debug().Model(&models.Building{}).Where("id=?", uid).Take(&models.Building{}).UpdateColumns(
			map[string]interface{}{
				"building_status": "0",
				"updated_at":      time.Now(),
			},
		)
		ch <- true
	}(done)

	if channels.OK(done) {
		if rs.Error != nil {
			return 0, rs.Error
		}
		return rs.RowsAffected, nil
	}
	return 0, rs.Error
}
func (r *repositoryBuildingCRUD) CheckUniqueBuilding(checkBuilding models.Building) (bool, error) {

	var err error
	buildings := models.Building{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.Building{}).Where("building_name=?", checkBuilding.Building_name).Take(&buildings).Error
		if err != nil {
			//record not found
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) { // found duplicate
		return true, errors.New("found duplicate building")
	} else {
		return false, nil

	}
}
