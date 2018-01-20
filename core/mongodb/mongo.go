package mongodb

import (
	"fmt"
	"github.com/dlintw/goconf"
	"gopkg.in/mgo.v2"
)

var (
	mgoSession *mgo.Session
	mgoUrl     string
)

func getSession() *mgo.Session {
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(mgoUrl)
		if err != nil {
			panic(err)
		}
	}
	return mgoSession.Clone()
}

func Mgo(db string, collection string, s func(*mgo.Collection) error) error {
	session := getSession()
	defer func() {
		session.Close()
		if err := recover(); err != nil {
			panic(err)
		}
	}()
	c := session.DB(db).C(collection)
	return s(c)
}

func InitMongoConnection(conf *goconf.ConfigFile) {
	mgoUrl, _ := conf.GetString("mongo", "mongo_url")
	if mgoUrl == "" {
		fmt.Println("未启用mongodb")
		return
	}
}
