package db

import "github.com/labstack/echo"

// WithDB middleware for Echo
func WithDB(db DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)

			return next(c)
		}
	}
}
