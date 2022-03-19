package crud

import (
	"apartment/api/models"
	"apartment/api/utils/channels"
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type repositoryUsersCRUD struct {
	db *gorm.DB
}

func NewRepositoryUsersCRUD(db *gorm.DB) *repositoryUsersCRUD {
	return &repositoryUsersCRUD{db}
}
func (r *repositoryUsersCRUD) FindAll() ([]models.User, error) {
	var err error
	users := []models.User{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		//Get only active user
		err = r.db.Debug().Model(&models.User{}).Where("User_status = '1' ").Limit(100).Find(&users).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return users, nil
	}
	return nil, err
}
func (r *repositoryUsersCRUD) FindByID(uid uint64) (models.User, error) {
	var err error
	user := models.User{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.User{}).Where("id = ?", uid).Take(&user).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return user, nil
	}
	if gorm.IsRecordNotFoundError(err) {
		return models.User{}, errors.New("user not found")
	}
	return models.User{}, err
}
func (r *repositoryUsersCRUD) Create(user models.User) (models.User, error) {
	var err error
	var found bool
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		//Check Unique    create user data not duplicate existing data (username , User_status=1 (active))
		found, err = r.CheckUniqueUser(user)
		if found { // found same data in db
			ch <- false
			return
		}
		err = r.db.Debug().Model(&models.User{}).Create(&user).Error
		if err != nil {
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) {
		return user, nil
	}
	return models.User{}, err
}
func (r *repositoryUsersCRUD) Update(uid uint64, user models.User) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		rs = r.db.Debug().Model(&models.User{}).Where("id = ?", uid).Take(&models.User{}).UpdateColumns(
			map[string]interface{}{
				"firstname":   user.Firstname,
				"lastname":    user.Lastname,
				"user_status": user.User_status,
				"user_level":  user.User_level,
				"updated_at":  time.Now(),
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
func (r *repositoryUsersCRUD) Delete(uid uint64) (int64, error) {
	var rs *gorm.DB
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)

		//Delete is update User_status = 0
		rs = r.db.Debug().Model(&models.User{}).Where("id=?", uid).Take(&models.User{}).UpdateColumns(
			map[string]interface{}{
				"user_status": "0",
				"updated_at":  time.Now(),
			},
		)
		//rs = r.db.Debug().Model(&models.User{}).Where("id = ?", uid).Take(&models.User{}).Delete(&models.User{})
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
func (r *repositoryUsersCRUD) CheckUniqueUser(checkUser models.User) (bool, error) {

	var err error
	user := models.User{}
	done := make(chan bool)
	go func(ch chan<- bool) {
		defer close(ch)
		err = r.db.Debug().Model(&models.User{}).Where("username=? and User_status=1", checkUser.Username).Take(&user).Error
		if err != nil {
			//record not found
			ch <- false
			return
		}
		ch <- true
	}(done)

	if channels.OK(done) { // found duplicate
		return true, errors.New("found duplicate user")
	} else {
		return false, nil

	}
}
