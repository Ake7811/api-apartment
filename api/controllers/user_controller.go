package controllers

import (
	"apartment/api/database"
	"apartment/api/models"
	"apartment/api/repository"
	"apartment/api/repository/crud"
	"encoding/json"
	"fmt"
	"strconv"

	//"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

func GetUsers(c echo.Context) error {
	db, err := database.Connect()
	if err != nil {
		c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer db.Close()
	repo := crud.NewRepositoryUsersCRUD(db)
	func(usersReposity repository.UserRepository) {
		users, err := usersReposity.FindAll()
		if err != nil {
			c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, users)
	}(repo)
	return nil
}

func GetUser(c echo.Context) error {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(c.JSON(http.StatusBadRequest, err.Error()))
		return nil
	}
	db, err := database.Connect()
	if err != nil {
		c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersCRUD(db)

	func(usersRepository repository.UserRepository) {
		user, err := usersRepository.FindByID(uint32(uid))
		if err != nil {
			c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, user)
	}(repo)
	return nil
}

func UpdateUser(c echo.Context) error {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(c.JSON(http.StatusBadRequest, err.Error()))
		return nil
	}

	//fmt.Printf("id : %d", uid)
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	//fmt.Printf("%#v\n%v\n", user, string(body))
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	db, err := database.Connect()
	if err != nil {
		c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersCRUD(db)

	func(usersRepository repository.UserRepository) {
		rows, err := usersRepository.Update(uint32(uid), user)
		if err != nil {
			c.Error(c.JSON(http.StatusBadRequest, err.Error()))
			return
		}
		c.JSON(http.StatusOK, rows)
	}(repo)
	return nil
}

func CreateUser(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	user := models.User{}
	err = json.Unmarshal(body, &user)
	//fmt.Printf("%#v\n%v\n", user, string(body))
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	user.Prepare()
	err = user.Validate("")
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	db, err := database.Connect()
	if err != nil {
		c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersCRUD(db)

	func(usersRepository repository.UserRepository) {
		user, err := usersRepository.Create(user)
		if err != nil {
			c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
			return
		}
		c.Response().Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().RequestURI, user.Id))
		c.JSON(http.StatusCreated, user)
	}(repo)

	return nil
}

func DeleteUser(c echo.Context) error {

	uid, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.Error(c.JSON(http.StatusBadRequest, err.Error()))
		return nil
	}
	db, err := database.Connect()
	if err != nil {
		c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer db.Close()

	repo := crud.NewRepositoryUsersCRUD(db)

	func(usersRepository repository.UserRepository) {
		_, err := usersRepository.Delete(uint32(uid))
		if err != nil {
			c.Error(c.JSON(http.StatusBadRequest, err.Error()))
			return
		}
		c.Response().Header().Set("Entity", fmt.Sprintf("%d", uid))
		c.JSON(http.StatusNoContent, "")

	}(repo)
	return nil
}
