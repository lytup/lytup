package models

import (
	"github.com/labstack/lytup/server/db"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type File struct {
	Id   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Size int64  `json:"size" bson:"size"`
	Type string `json:"type" bson:"type"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func FindFileById(folderId, fileId string) *File {
	db := db.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	err := db.Collection.Find(bson.M{"id": folderId}).Select(bson.M{"files": bson.M{"$elemMatch": bson.M{"id": fileId}}}).One(&fol)
	if err != nil {
		log.Fatal(err)
	}
	return fol.Files[0]
}
