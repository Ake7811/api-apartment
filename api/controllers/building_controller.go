package controllers

import (
	"apartment/api/auth"
	"apartment/api/database"
	"apartment/api/models"
	"apartment/api/repository"
	"apartment/api/repository/crud"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

func GetBuildings(c echo.Context) error {
	db, err := database.Connect()
	if err != nil {
		c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
		return nil
	}
	defer db.Close()
	repo := crud.NewRepositoryBuildingCRUD(db)
	func(buildingsReposity repository.BuildingRepository) {
		users, err := buildingsReposity.FindAll()
		if err != nil {
			c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, users)
	}(repo)
	return nil
}

func GetBuilding(c echo.Context) error {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
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

	repo := crud.NewRepositoryBuildingCRUD(db)

	func(buildingsRepository repository.BuildingRepository) {
		user, err := buildingsRepository.FindByID(uint64(uid))
		if err != nil {
			c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, user)
	}(repo)
	return nil
}

func UpdateBuilding(c echo.Context) error {
	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
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
	building := models.Building{}
	err = json.Unmarshal(body, &building)
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

	repo := crud.NewRepositoryBuildingCRUD(db)

	func(buildingsRepository repository.BuildingRepository) {
		rows, err := buildingsRepository.Update(uint64(uid), building)
		if err != nil {
			c.Error(c.JSON(http.StatusBadRequest, err.Error()))
			return
		}
		c.JSON(http.StatusOK, rows)
	}(repo)
	return nil
}

func CreateBuilding(c echo.Context) error {

	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}
	building := models.Building{}
	err = json.Unmarshal(body, &building)
	if err != nil {
		c.Error(c.JSON(http.StatusUnprocessableEntity, err.Error()))
		return nil
	}

	tokenAuth, err := auth.ExtractTokenMetadata(c.Request())
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return nil
	}

	userId, err := auth.FetchAuth(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return nil
	}
	building.Create_by = userId

	building.Prepare()
	err = building.Validate("")
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

	repo := crud.NewRepositoryBuildingCRUD(db)

	func(buildingsRepository repository.BuildingRepository) {
		building, err := buildingsRepository.Create(building)
		if err != nil {
			c.Error(c.JSON(http.StatusInternalServerError, err.Error()))
			return
		}
		c.Response().Header().Set("Location", fmt.Sprintf("%s%s/%d", c.Request().Host, c.Request().RequestURI, building.Id))
		c.JSON(http.StatusCreated, building)
	}(repo)

	return nil
}

func DeleteBuilding(c echo.Context) error {

	uid, err := strconv.ParseUint(c.Param("id"), 10, 64)
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

	repo := crud.NewRepositoryBuildingCRUD(db)

	func(buildingsRepository repository.BuildingRepository) {
		_, err := buildingsRepository.Delete(uint64(uid))
		if err != nil {
			c.Error(c.JSON(http.StatusBadRequest, err.Error()))
			return
		}
		c.Response().Header().Set("Entity", fmt.Sprintf("%d", uid))
		c.JSON(http.StatusNoContent, "")

	}(repo)
	return nil
}
