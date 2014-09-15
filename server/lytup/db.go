package lytup

import (
	"fmt"

	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

var (
	session *mgo.Session
)

type Db struct {
	Session    *mgo.Session
	Collection *mgo.Collection
}

func (db *Db) Close() {
	db.Session.Close()
}

func NewDb(collection string) *Db {
	return &Db{session.Copy(), session.DB("lytup").C(collection)}
}

func init() {
	uri := "mongodb://"

	if C.MongoDb.Username != "" && C.MongoDb.Password != "" {
		uri = fmt.Sprintf("%s%s:%s@", uri, C.MongoDb.Username, C.MongoDb.Password)
	}
	uri += C.MongoDb.Host
	if C.MongoDb.Port != 0 {
		uri = fmt.Sprintf("%s:%d", uri, C.MongoDb.Port)
	}

	var err error
	if session, err = mgo.Dial(uri); err != nil {
		glog.Fatal(err)
	}

	// Set indexes
	db := NewDb("users")
	if err := db.Collection.EnsureIndex(mgo.Index{
		Key:    []string{"email"},
		Unique: true,
	}); err != nil {
		glog.Fatal(err)
	}
}
