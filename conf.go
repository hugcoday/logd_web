package main

import "logd_web/app/utils"

// Conf 配置入口
func Conf() {
	// cfg, err := ini.LooseLoad("./conf/config.ini", "filename_404")

	// if err != nil {
	// 	fmt.Println("conf properties is not found")
	// }

	// mysql, err := cfg.GetSection("mysql")
	// if err != nil {
	// 	log.Error(err)
	// }
	// dbURL := mysql.Key("db_url").Value()

	// fmt.Println(dbURL)

	cfg := utils.ReadConfig("./conf/config.ini")

	Initlog(cfg["log"])

}
