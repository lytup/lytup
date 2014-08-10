package db

import (
	"github.com/golang/glog"
	"labix.org/v2/mgo"
)

var session *mgo.Session

type Db struct {
	Session    *mgo.Session
	Collection *mgo.Collection
}

func init() {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		glog.Fatal(err)
	}
	session = sess
}

func NewDb(collection string) *Db {
	return &Db{session.Copy(), session.DB("lytup").C(collection)}
}
