package api

import (
	"apartment/api/routes"
	"apartment/config"
	"fmt"

	"github.com/labstack/echo"
)

//

func Run() {
	config.Load()

	fmt.Printf("running... at port %d", config.PORT)

	e := echo.New()
	routes.InitUserRoutes(e)
	routes.InitLoginRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.PORT)))
}
