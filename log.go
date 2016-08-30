package main

import (
	log "github.com/gogap/logrus"
	"github.com/gogap/logrus/hooks/file"
)

var Log *log.Logger

func Initlog(config map[string]string) {
	log.SetFormatter(&log.TextFormatter{})

	//输出到graylog
	// glog, err := graylog.NewHook("boot2docker:9001", "yijifu", nil)
	// if err != nil {
	// 	log.Error(err)
	// 	return
	// }
	// log.AddHook(glog)

	//输出到文件
	log.AddHook(file.NewHook(config["log_file"]))

	//yijifu组件中的member模块的日志
	//log.WithField("biz", "member").Errorf("member not login,member is %s", "1001")

}
