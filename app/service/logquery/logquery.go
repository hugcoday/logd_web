package logquery

import (
	"net/http"

	"github.com/labstack/echo"
)

// Index  获取用户列表
func Index(c echo.Context) error {

	return c.String(http.StatusOK, "no auth!")
}

func acountList(c echo.Context) error {

	return c.String(http.StatusOK, "no auth!")
}

// Route 路由
func Route(e *echo.Echo) {
	e.GET("/login", Index)
	e.POST("/user/account", acountList)

}
