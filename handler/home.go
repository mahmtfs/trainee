package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func HomeHandler(c echo.Context) error{
	return c.Render(http.StatusOK, "index.html", nil)
}
