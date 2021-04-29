package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *handler) Home(c echo.Context) error {
	return c.String(http.StatusOK, "welcome home")
}
