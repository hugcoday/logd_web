package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"

	log "github.com/gogap/logrus"
)

func main() {
	//config init
	initlog()
	Conf()

	e := echo.New()

	//middleware config init
	//route log
	//	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	//		Format: "time=${time_rfc3339},method=${method}, uri=${uri}, status=${status}\n",
	//	}))
	// gzip open
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	// static file

	e.Use(middleware.CORS())

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   "public",
		Browse: false,
	}))

	//	e.Use(middleware.JWTWithConfig(middleware.JWTConfig{
	//		SigningKey: []byte("secret"),
	//		Extractor:  middleware.JWTFromHeader,
	//	}))
	// route init
	Route(e)
	log.Info("server started at 3002")
	e.Run(fasthttp.New(":3002"))

}
