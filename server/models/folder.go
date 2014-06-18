package models

import (
	"github.com/dchest/uniuri"
	"github.com/labstack/lytup/server/db"
	"labix.org/v2/mgo/bson"
	"time"
)

// TODO: move status to enum like structure
type Folder struct {
	Id        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Files     []*File   `json:"files" bson:"files"`
	Status    string    `json:"status" bson:"status"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	ExpiresAt time.Time `json:"expiresAt" bson:"expiresAt"`
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

func (fol *Folder) Save() {
	fol.Id = uniuri.NewLen(5)
	fol.Files = []*File{}
	fol.CreatedAt = time.Now()
	fol.ExpiresAt = fol.CreatedAt.Add(24 * time.Hour)

	db := db.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Insert(&fol)
	if err != nil {
		panic(err)
	}
}

func UpdateFolder(id string, fol *Folder) {
	db := db.NewDb("folders")
	defer db.Session.Close()

	fol.UpdatedAt = time.Now()

	// Files
	if fol.Files != nil {
		for _, file := range fol.Files {
			file.Id = "uniuri.NewLen(5)"
			file.CreatedAt = time.Now()
		}

		err := db.Collection.Update(bson.M{"id": id}, bson.M{"$set": bson.M{"updatedAt": time.Now()}, "$pushAll": bson.M{"files": fol.Files}})
		if err != nil {
			panic(err)
		}
	}

	// TODO: Update other fields on-demand basis
}
