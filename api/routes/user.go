package routes

import (
	"apartment/api/controllers"

	"github.com/labstack/echo"
)

func InitUserRoutes(e *echo.Echo) {
	e.GET("/users", controllers.GetAllUser)
}
