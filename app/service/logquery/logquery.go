package logquery

import (
	"logd_web/app/model"
	"logd_web/comm/utils"
	"logd_web/comm/validate"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	log "github.com/gogap/logrus"
	"github.com/labstack/echo"
)

// SearchData : query condition
// type SearchData struct {
// 	QueryDate int    `json:"querydate"`
// 	QueryText string `json:"querytext"`
// }

// Index  获取用户列表
func Index(c echo.Context) error {
	log.Info("start")
	//var result LogData
	condition := new(SearchData)
	if err := c.Bind(condition); err != nil {
		log.Error(err)
		return err
	}
	// 参数错误
	errs := validate.Vd.Struct(condition)

	if errs != nil {
		log.Error(errs)
		result := model.ResultModel{1001, "请求参数错误", nil}
		return c.JSON(http.StatusOK, result)
	}

	result := model.ResultModel{1000, "ok", nil}
	pages := new(model.PagesModel)
	if err := c.Bind(pages); err != nil {
		return err
	}
	curIndex := 0
	if pages.PageIndex > 1 {
		curIndex = pages.PageIndex - 1
	}
	// pages := new(model.PagesModel)
	// if err := c.Bind(pages); err != nil {
	// 	return err
	// }

	//	mmap := utils.ConvertToMap(condition)
	var queryCondition bson.M
	if len(condition.QueryText) == 0 {
		queryCondition = bson.M{"ctime": bson.M{"$gte": condition.StartDate, "$lte": condition.EndDate}}
	} else {
		queryCondition = bson.M{"ctime": bson.M{"$gte": condition.StartDate, "$lte": condition.EndDate}, "$text": bson.M{"$search": condition.QueryText}}
	}

	data := utils.MDB.FindPage("logd_20160829", "logdata", queryCondition, curIndex)
	count := utils.MDB.FindCount("logd_20160829", "logdata", queryCondition)
	//pages := base.PagesModel{}
	pages.Total = count
	pages.Data = data

	pages.PageSize = 10
	result.Extra = pages
	log.Info(result)
	return c.JSON(http.StatusOK, result)
}

// Route 路由
func Route(e *echo.Echo) {
	e.POST("/log/query", Index)

}
