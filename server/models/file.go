package models

import (
	"github.com/dchest/uniuri"
	"github.com/labstack/lytup/server/db"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"os"
	"path"
	"time"
)

type File struct {
	Id        string    `json:"id,omitempty" bson:"id"`
	Name      string    `json:"name,omitempty" bson:"name"`
	Size      uint64    `json:"size,omitempty" bson:"size"`
	Type      string    `json:"type,omitempty" bson:"type"`
	Loaded    uint8     `json:"loaded,omitempty" bson:"loaded"`
	Uri       string    `json:"uri,omitempty" bson:"uri"`
	Thumbnail string    `json:"thumbnail,omitempty" bson:"thumbnail"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt"`
}

func (file *File) Create(folId string) {
	file.Id = uniuri.NewLen(7)
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

func FindFileById(id string) (string, *File) {
	db := db.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	err := db.Collection.Find(bson.M{"files.id": id}).
		Select(bson.M{"id": 1, "files": bson.M{"$elemMatch": bson.M{"id": id}}}).
		One(&fol)
	if err != nil {
		panic(err)
	}

	return fol.Id, fol.Files[0]
}

func (file *File) Update(folId string) {
	now := time.Now()
	m := bson.M{"updatedAt": now}

	if file.Loaded != 0 {
		m["files.$.loaded"] = file.Loaded
	}

	if file.Uri != "" {
		m["files.$.uri"] = file.Uri
	}

	if file.Thumbnail != "" {
		m["files.$.thumbnail"] = file.Thumbnail
	}

	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Update(bson.M{"id": folId, "files.id": file.Id},
		bson.M{"$set": m})
	if err != nil {
		panic(err)
	}
}

func DeleteFile(folId, fileId string) {
	db := db.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	_, err := db.Collection.Find(bson.M{"id": folId, "files.id": fileId}).
		Select(bson.M{"files": bson.M{"$elemMatch": bson.M{"id": fileId}}}).
		Apply(mgo.Change{Update: bson.M{"$pull": bson.M{"files": bson.M{"id": fileId}}}}, &fol)
	if err != nil {
		panic(err)
	}

	file := fol.Files[0]

	// Delete file from file system
	err = os.Remove(path.Join("/tmp", folId, file.Name)) // TODO: Read from config
	if err != nil {
		panic(err)
	}
}
