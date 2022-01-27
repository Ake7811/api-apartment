package controllers

import (
	"apartment/api/auth"
	"apartment/api/models"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

func Login(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}

	user := models.User{}
	err = json.Unmarshal(body, &user)
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	user.Prepare()
	err = user.Validate("login")
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	token, err := auth.SignIn(user.Username, user.Password)
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	return c.JSON(http.StatusOK, token)

	



}
