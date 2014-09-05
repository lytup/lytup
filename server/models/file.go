package models

import (
	"os"
	"path"
	"time"

	"github.com/dchest/uniuri"
	L "github.com/labstack/lytup/server/lytup"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (file *File) Create(folId string, usr *User) error {
	file.Id = uniuri.NewLen(7)
	file.CreatedAt = time.Now()

	db := L.NewDb("folders")
	defer db.Session.Close()
	if err := db.Collection.Update(bson.M{"id": folId, "userId": usr.Id},
		bson.M{"$set": bson.M{"updatedAt": time.Now()},
			"$push": bson.M{"files": file}}); err != nil {
		return err
	}
	return nil
}

func FindFileById(id string) (*File, string, error) {
	db := L.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	if err := db.Collection.Find(bson.M{"files.id": id}).
		Select(bson.M{"id": 1, "files": bson.M{"$elemMatch": bson.M{"id": id}}}).
		One(&fol); err != nil {
		return nil, "", err
	}
	return fol.Files[0], fol.Id, nil
}

func (file *File) Update(folId string, usr *User) error {
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

	db := L.NewDb("folders")
	defer db.Session.Close()
	if err := db.Collection.Update(bson.M{"id": folId, "userId": usr.Id,
		"files.id": file.Id}, bson.M{"$set": m}); err != nil {
		return err
	}
	return nil
}

func DeleteFile(folId, fileId string, usr *User) error {
	db := L.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	if _, err := db.Collection.Find(bson.M{"id": folId, "userId": usr.Id, "files.id": fileId}).
		Select(bson.M{"files": bson.M{"$elemMatch": bson.M{"id": fileId}}}).
		Apply(mgo.Change{Update: bson.M{"$pull": bson.M{"files": bson.M{"id": fileId}}}}, &fol); err != nil {
		return err
	}

	file := fol.Files[0]

	// Delete file from file system
	if err := os.Remove(path.Join("/tmp", folId, file.Name)); err != nil { // TODO: Read from config
		return err
	}

	return nil
}
