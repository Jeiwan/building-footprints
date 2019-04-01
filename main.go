package main

import (
	"fmt"
	"os"

	"github.com/Jeiwan/building-footprints/handlers"
	"github.com/go-playground/validator"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	app := cli.NewApp()
	app.Name = "building-footprints"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port",
			Value: 3000,
		},
	}
	app.Action = func(c *cli.Context) error {

		e := echo.New()
		e.Use(
			middleware.Logger(),
			middleware.Recover(),
			middleware.Gzip(),
		)
		e.Validator = &customValidator{validator: validator.New()}

		api := e.Group("/api/v0")
		api.GET("/avg_area", handlers.AvgArea)

		return e.Start(fmt.Sprintf(":%d", c.Int("port")))
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln(err)
	}
}
