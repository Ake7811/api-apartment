package routes

import (
	"apartment/api/controllers"

	"github.com/labstack/echo"
)

func InitUserRoutes(e *echo.Echo) {
	e.GET("/users", controllers.GetUsers)
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/:id", controllers.GetUser)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)
}

func InitLoginRoutes(e *echo.Echo) {
	e.POST("/login", controllers.Login)
}
