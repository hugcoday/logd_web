package logquery

import (
	"net/http"

	"estate/app/model"

	"estate/comm/utils"
	"estate/comm/validate"

	"gopkg.in/mgo.v2/bson"

	log "github.com/hugcoday/logrus"
	"github.com/kataras/iris"
)

// Index  获取用户列表
func Index(c iris.Context) error {
	//var result LogData
	condition := new(SearchData)
	if err := c.ReadJSON(condition); err != nil {
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
	if err := c.ReadJSON(pages); err != nil {
		return err
	}
	curIndex := 0
	if pages.PageIndex > 1 {
		curIndex = pages.PageIndex - 1
	}

	var queryCondition bson.M
	if len(condition.QueryText) == 0 {
		queryCondition = bson.M{"ctime": bson.M{"$gte": condition.StartDate, "$lte": condition.EndDate}}
	} else {
		queryCondition = bson.M{"ctime": bson.M{"$gte": condition.StartDate, "$lte": condition.EndDate}, "$text": bson.M{"$search": condition.QueryText}}
	}

	data, count := utils.MDB.FindPage("logd_20160829", "logdata", queryCondition, curIndex)

	pages.Total = count
	pages.Data = data

	pages.PageSize = 10
	result.Extra = pages
	log.Info(result)
	return c.JSON(http.StatusOK, result)
}

// Route 路由
func Route() {
	iris.Post("/log/query", Index)
}
