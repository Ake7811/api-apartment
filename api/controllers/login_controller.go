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

	//Login by user and password  got Set of token (Access token , Refresh token)
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

func Logout(c echo.Context) error {

	result, err := auth.SignOut(c)
	if err != nil {
		c.Error(c.JSON(http.StatusUnauthorized, err.Error()))
		return nil
	}

	if result {
		return c.JSON(http.StatusOK, "Successfully logged out")
	} else {
		return c.JSON(http.StatusBadRequest, nil)
	}
}
