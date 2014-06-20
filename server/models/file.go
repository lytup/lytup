package models

import (
	"github.com/dchest/uniuri"
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
	Loaded int `json:"loaded" bson:"loaded"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func (file *File) Create(folderId string) {
	file.Id = uniuri.NewLen(5)
	file.CreatedAt = time.Now()

	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Update(bson.M{"id": folderId},
			bson.M{"$set": bson.M{"updatedAt": time.Now()},
				"$push": bson.M{"files": file}})
	if err != nil {
		panic(err)
	}
}

func UpdateFile(folderId, fileId string, file *File) {
	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Update(bson.M{"id": folderId, "files.id": fileId},
			bson.M{"$set": bson.M{"updatedAt": time.Now(),
				"files.$.loaded": file.Loaded}})
	if err != nil {
		panic(err)
	}
}

func FindFileById(folderId, fileId string) *File {
	db := db.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	err := db.Collection.Find(bson.M{"id": folderId}).
			Select(bson.M{"files": bson.M{"$elemMatch": bson.M{"id": fileId}}}).
			One(&fol)
	if err != nil {
		log.Fatal(err)
	}
	return fol.Files[0]
}
