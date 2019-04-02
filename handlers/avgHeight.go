package handlers

import (
	"net/http"

	"github.com/labstack/echo"
)

type request struct {
	BoroughCode int `query:"borough_code" validate:"required"`
}

type response float64

// AvgHeight renders average height by borough code
func AvgHeight(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "borough_code is not specified"})
	}

	var resp response

	return c.JSON(http.StatusOK, resp)
}
