package main

import (
	"logd_web/app/service/logquery"

	"github.com/labstack/echo"
)

// Route init
func Route(e *echo.Echo) {

	logquery.Route(e)
}
