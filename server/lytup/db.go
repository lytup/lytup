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
	cfg := Config.MongoDb
	uri := "mongodb://"

	if cfg.Username != "" && cfg.Password != "" {
		uri = fmt.Sprintf("%s%s:%s@", uri, cfg.Username, cfg.Password)
	}
	uri += cfg.Host
	if cfg.Port != 0 {
		uri = fmt.Sprintf("%s:%d", uri, cfg.Port)
	}

	var err error
	if session, err = mgo.Dial(uri); err != nil {
		glog.Fatal(err)
	}
}

func NewDb(collection string) *Db {
	return &Db{session.Copy(), session.DB("lytup").C(collection)}
}
