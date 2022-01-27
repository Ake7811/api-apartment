package auth

import (
	"apartment/api/database"
	"apartment/api/models"
	"apartment/api/security"
	"apartment/api/utils/channels"

	"github.com/jinzhu/gorm"
)

func SignIn(user, password string) (string, error) {
	userQuery := models.User{}
	var err error
	var db *gorm.DB
	done := make(chan bool)

	go func(ch chan<- bool) {
		defer close(ch)
		db, err = database.Connect()
		if err != nil {
			ch <- false
			return
		}
		defer db.Close()

		err = db.Debug().Model(models.User{}).Where("username = ?", user).Take(&userQuery).Error
		if err != nil {
			ch <- false
			return
		}

		err = security.VerifyPassword(userQuery.Password, password)
		if err != nil {
			ch <- false
			return

		}
		ch <- true

	}(done)

	if channels.OK(done) {
		return CreateToken(userQuery.Id)
	}
	return "", err

}
