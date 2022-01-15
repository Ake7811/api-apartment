package controllers

import (
	"apartment/api/database"
	"apartment/api/repository"
	"apartment/api/repository/crud"
	"net/http"

	"github.com/labstack/echo"
)

func GetAllUser(c echo.Context) error {
	db, err := database.Connect()
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	defer db.Close()
	repo := crud.NewRepositoryUsersCRUD(db)
	func(usersReposity repository.UserRepository) {
		users, err := usersReposity.FindAll()
		if err != nil {
			c.NoContent(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, users)
	}(repo)
	return nil
}
