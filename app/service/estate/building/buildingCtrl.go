package building

import (
	"github.com/kataras/iris"

	log "github.com/hugcoday/logrus"
)

//Route init()
func Route() {
	iris.Post("/building/query", Query)
}

//Query 查询
func Query(ctx *iris.Context) {
	ctx.JSON(iris.StatusOK, "text")
	log.Info("query test")
}
