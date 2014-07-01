package models

import (
	"github.com/dchest/uniuri"
	"github.com/labstack/lytup/server/db"
	"labix.org/v2/mgo/bson"
	"os"
	"path"
	"time"
)

type Folder struct {
	Id        string        `json:"id" bson:"id"`
	Name      string        `json:"name" bson:"name"`
	Files     []*File       `json:"files" bson:"files"`
	UserId    bson.ObjectId `json:"userId" bson:"userId"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
	ExpiresAt int64         `json:"expiresAt" bson:"expiresAt"`
}

func (fol *Folder) Create() {
	fol.Id = uniuri.NewLen(5)
	fol.Files = []*File{}
	fol.CreatedAt = time.Now()
	fol.ExpiresAt = fol.CreatedAt.Add(4 * time.Hour).Unix()

	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Insert(&fol)
	if err != nil {
		panic(err)
	}
}

func FindFolders() *[]Folder {
	db := db.NewDb("folders")
	defer db.Session.Close()
	folders := []Folder{}
	err := db.Collection.Find(nil).All(&folders)
	if err != nil {
		panic(err)
	}
	return &folders
}

func FindFolderById(id string) *Folder {
	db := db.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	err := db.Collection.Find(bson.M{"id": id}).One(&fol)
	if err != nil {
		panic(err)
	}
	return &fol
}

func UpdateFolder(id string, fol *Folder) {
	fol.UpdatedAt = time.Now()

	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Update(bson.M{"id": id},
		bson.M{"$set": bson.M{"updatedAt": time.Now()}})
	if err != nil {
		panic(err)
	}

	// TODO: Update other fields on-demand basis
}

func DeleteFolder(id string) {
	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Remove(bson.M{"id": id})
	if err != nil {
		panic(err)
	}

	// Delete folder from file system
	err = os.RemoveAll(path.Join("/tmp", id)) // TODO: Read from config
	if err != nil {
		panic(err)
	}
}
