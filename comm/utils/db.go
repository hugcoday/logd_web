package utils

import (
	"bytes"
	"fmt"
	"logd/lib"

	"os"
	"time"

	log "github.com/gogap/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Mdb struct {
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

var MDB Mdb

func InitDB(config map[string]string) (mo Mdb) {

	mo.wq = lib.NewWaitQuit("mongodb outputer", -1)
	mo.mongosAddr, _ = config["mongos"]
	mo.db, _ = config["db"]
	mo.collection, _ = config["collection"]
	mo.session = initMongoDbSession(mo.mongosAddr)
	//暂时出错直接退出
	if mo.session == nil {
		log.Error("init mongodb session failed")
		os.Exit(1)
	}

	mo.transactionIdKey = "transaction_id"
	mo.fileList = lib.GlobalListInit()
	MDB = mo
	return mo
}

func GetDB() (mo Mdb) {
	return MDB
}

func initMongoDbSession(mongosAddr string) *mgo.Session {
	session, err := mgo.Dial(mongosAddr)
	if err != nil {
		log.Error(fmt.Sprintf("init mongodb session error:%v", err))
		return nil
	}

	session.SetMode(mgo.Monotonic, true) //设置read preference
	session.SetSafe(&mgo.Safe{W: 1})     //设置write concern
	return session
}

//用于检验session可用性并适时重连的routine func
//用于重连main session
func (this *Mdb) reConnMongoDb() {
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
				log.Info("session re-dial failed!")
			} else {
				log.Info("session re-dial success!")
			}
		}
		time.Sleep(time.Second)
	}
}

//用于routine重新clone session, main session未重连，则继续用旧的session
func (this *Mdb) reCloneRoutineSession(psession **mgo.Session) {
	if this.session != nil {
		//re-clone a session
		*psession = this.session.Clone()
	}
}

func (this *Mdb) Find(db string, collection string, condition bson.M) (results []interface{}) {
	//	var session *mgo.Session
	c := this.session.DB(db).C(collection)

	err := c.Find(condition).Sort("-ctime").All(&results)
	if err != nil {
		log.Error(err)
	}
	return results
}

// FindCount dsd
func (this *Mdb) FindCount(db string, collection string, condition bson.M) (result int) {
	//	var session *mgo.Session
	c := this.session.DB(db).C(collection)

	result, err := c.Find(condition).Count()
	if err != nil {
		log.Error(err)
	}
	return result
}

func (this *Mdb) FindPage(db string, collection string, condition bson.M, PageIndex int) (results []interface{}) {
	//	var session *mgo.Session
	c := this.session.DB(db).C(collection)

	err := c.Find(condition).Sort("-ctime").Skip(PageIndex * 10).Limit(10).All(&results)
	if err != nil {
		log.Error(err)
	}
	return results
}
