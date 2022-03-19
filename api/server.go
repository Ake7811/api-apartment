package api

import (
	"apartment/api/auth"
	"apartment/api/controllers"
	"apartment/api/middleware"
	"apartment/config"
	"fmt"

	"github.com/labstack/echo"
)

//

func Run() {
	config.Load()

	fmt.Printf("running... at port %d", config.PORT)

	e := echo.New()
	e.POST("/login", controllers.Login)
	e.POST("/token/refresh", auth.TokenRefresh)

	auth := e.Group("/auth", middleware.TokenAuthMiddleware)
	auth.POST("/logout", controllers.Logout)

	//User
	auth.GET("/users", controllers.GetUsers)
	auth.POST("/users", controllers.CreateUser)
	auth.GET("/users/:id", controllers.GetUser)
	auth.PUT("/users/:id", controllers.UpdateUser)
	auth.DELETE("/users/:id", controllers.DeleteUser)

	//Building
	auth.GET("/buildings", controllers.GetBuildings)
	auth.POST("/buildings", controllers.CreateBuilding)
	auth.GET("/buildings/:id", controllers.GetBuilding)
	auth.PUT("/buildings/:id", controllers.UpdateBuilding)
	auth.DELETE("/buildings/:id", controllers.DeleteBuilding)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))
}
