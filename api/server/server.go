package server

import (
	"api/server/route"
	"context"

	"github.com/labstack/echo/v4"
)

func Run(ctx context.Context) error {
	e := echo.New()
	route.InitRoute(e)

	return e.Start(":8080")
}
