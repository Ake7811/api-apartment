package middleware

import (
	"apartment/api/auth"
	"net/http"

	"github.com/labstack/echo"
)

func TokenAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := auth.TokenValid(c.Request())
		if err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			return nil
		}

		return next(c)
	}
}
