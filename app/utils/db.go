package utils

import (
	"bytes"
	"fmt"
	"logd/lib"
	"logd/loglib"
	"os"
	"time"

	"gopkg.in/mgo.v2"
)

type MongoDbOutputer struct {
	buffer               chan bytes.Buffer
	mongosAddr           string
	session              *mgo.Session
	db                   string
	collection           string
	isUpsert             bool
	bulkSize             int
	savers               int
	file_mem_folder_name string
	transactionIdKey     string
	fileList             *lib.GlobalList
	wq                   *lib.WaitQuit
}

func InitDB(config map[string]string) (mo MongoDbOutputer) {

	mo.wq = lib.NewWaitQuit("mongodb outputer", -1)
	mo.mongosAddr, _ = config["mongos"]
	mo.db, _ = config["db"]
	mo.collection, _ = config["collection"]
	mo.session = initMongoDbSession(mo.mongosAddr)
	//暂时出错直接退出
	if mo.session == nil {
		loglib.Error("init mongodb session failed")
		os.Exit(1)
	}

	mo.transactionIdKey = "transaction_id"
	mo.fileList = lib.GlobalListInit()
	return mo
}

func initMongoDbSession(mongosAddr string) *mgo.Session {
	session, err := mgo.Dial(mongosAddr)
	if err != nil {
		loglib.Error(fmt.Sprintf("init mongodb session error:%v", err))
		return nil
	}

	session.SetMode(mgo.Monotonic, true) //设置read preference
	session.SetSafe(&mgo.Safe{W: 1})     //设置write concern
	return session
}

//用于检验session可用性并适时重连的routine func
//用于重连main session
func (this *MongoDbOutputer) reConnMongoDb() {
	nPingFail := 0 //ping失败次数
	reDial := false
	for {
		reDial = false
		if this.session == nil {
			//session未初始化
			reDial = true
		} else if this.session.Ping() != nil {
			//session连接丢失
			nPingFail++
			if nPingFail == 3 {
				reDial = true
			}
		}

		if reDial {
			nPingFail = 0
			this.session = initMongoDbSession(this.mongosAddr)
			if this.session == nil {
				loglib.Info("session re-dial failed!")
			} else {
				loglib.Info("session re-dial success!")
			}
		}
		time.Sleep(time.Second)
	}
}

//用于routine重新clone session, main session未重连，则继续用旧的session
func (this *MongoDbOutputer) reCloneRoutineSession(psession **mgo.Session) {
	if this.session != nil {
		//re-clone a session
		*psession = this.session.Clone()
	}
}
