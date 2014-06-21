package models

import (
	"github.com/dchest/uniuri"
	"github.com/labstack/lytup/server/db"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

type File struct {
	Id        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Size      int64     `json:"size" bson:"size"`
	Type      string    `json:"type" bson:"type"`
	Loaded    int       `json:"loaded" bson:"loaded"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func (file *File) Create(folId string) {
	file.Id = uniuri.NewLen(5)
	file.CreatedAt = time.Now()

	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Update(bson.M{"id": folId},
		bson.M{"$set": bson.M{"updatedAt": time.Now()},
			"$push": bson.M{"files": file}})
	if err != nil {
		panic(err)
	}
}

func UpdateFile(folId, id string, file *File) {
	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Update(bson.M{"id": folId, "files.id": id},
		bson.M{"$set": bson.M{"updatedAt": time.Now(),
			"files.$.loaded": file.Loaded}})
	if err != nil {
		panic(err)
	}
}

func FindFileById(id string) (string, *File) {
	db := db.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	err := db.Collection.Find(bson.M{"files.id": id}).
		Select(bson.M{"id": 1, "files.$": 1}).
		One(&fol)
	if err != nil {
		log.Fatal(err)
	}
	return fol.Id, fol.Files[0]
}
