package db

import (
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
		panic(err)
	}
	session = sess
}

func NewDb(collection string) *Db {
	db := Db{session.Copy(), session.DB("lytup").C(collection)}
	return &db
}
