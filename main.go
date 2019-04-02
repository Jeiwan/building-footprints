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
		cli.StringFlag{
			Name:  "mongo-url",
			Value: "127.0.0.1:27017",
		},
		cli.StringFlag{
			Name:  "mongo-db-name",
			Value: "building-footprints",
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
		api.GET("/avg_height", handlers.AvgHeight)

		return e.Start(fmt.Sprintf(":%d", c.Int("port")))
	}

	app.Commands = []cli.Command{
		cli.Command{
			Name: "load-data",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "mongo-url",
					Value: "127.0.0.1:27017",
				},
				cli.StringFlag{
					Name:  "mongo-db-name",
					Value: "building-footprints",
				},
				cli.StringFlag{
					Name:  "data-file",
					Value: "rows.json",
				},
			},
			Action: cliLoadData,
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatalln(err)
	}
}
