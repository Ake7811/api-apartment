package auth

import (
	"apartment/api/database"
	"apartment/api/models"
	"apartment/api/security"
	"apartment/api/utils/channels"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func SignIn(username, password string) (map[string]string, error) {
	user := models.User{}
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

		err = db.Debug().Model(models.User{}).Where("username = ?", username).Take(&user).Error
		if err != nil {
			ch <- false
			return
		}

		err = security.VerifyPassword(user.Password, password)
		if err != nil {
			ch <- false
			return

		}
		ch <- true

	}(done)

	if channels.OK(done) {
		ts, err := CreateToken(user.Id)
		if err != nil {
			return nil, err
		}

		saveErr := CreateAuth(user.Id, ts)
		if saveErr != nil {
			return nil, saveErr
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		return tokens, saveErr
	}

	return nil, err
}

func SignOut(c echo.Context) (bool, error) {
	au, err := ExtractTokenMetadata(c.Request())
	if err != nil {
		return false, err
	}
	deleted, delErr := DeleteAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		return false, delErr
	}
	return true, nil
}
