package handlers

import (
	"net/http"

	"github.com/Jeiwan/building-footprints/db"
	"github.com/labstack/echo"
)

type request struct {
	BoroughCode int `query:"borough_code" validate:"required"`
}

type response struct {
	AvgHeigh float64 `json:"avg_height"`
}

// AvgHeight renders average height by borough code
func AvgHeight(c echo.Context) error {
	var req request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "borough_code is not specified"})
	}

	db := c.Get("db").(db.DB)
	var resp response
	avgHeight, err := db.AvgHeightByBoroughCode(req.BoroughCode)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	resp.AvgHeigh = avgHeight

	return c.JSON(http.StatusOK, resp)
}
