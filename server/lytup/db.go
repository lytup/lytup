package lytup

import (
	"fmt"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

var session *mgo.Session

type Db struct {
	Session    *mgo.Session
	Collection *mgo.Collection
}

func init() {
	m := Config.MongoDb
	uri := "mongodb://"

	if m.Username != "" && m.Password != "" {
		uri = fmt.Sprintf("%s%s:%s@", uri, m.Username, m.Password)
	}
	uri += m.Host
	if m.Port != 0 {
		uri = fmt.Sprintf("%s:%d", uri, m.Port)
	}

	var err error
	if session, err = mgo.Dial(uri); err != nil {
		glog.Fatal(err)
	}
}

func NewDb(collection string) *Db {
	return &Db{session.Copy(), session.DB("lytup").C(collection)}
}
