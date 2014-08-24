package models

import (
	L "github.com/labstack/lytup/server/lytup"
	U "github.com/labstack/lytup/server/utils"
	"gopkg.in/mgo.v2/bson"
	"os"
	"path"
	"time"
)

type Folder struct {
	Id        string        `json:"id" bson:"id"`
	Name      string        `json:"name" bson:"name"`
	Files     []*File       `json:"files" bson:"files"`
	Expiry    uint16        `json:"expiry" bson:"expiry"`
	UserId    bson.ObjectId `json:"userId" bson:"userId"`
	CreatedAt time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt" bson:"updatedAt"`
	ExpiresAt time.Time     `json:"expiresAt" bson:"expiresAt"`
}

func (fol *Folder) Create() {
	fol.Id = U.RandString(7)
	fol.Files = []*File{}
	fol.CreatedAt = time.Now()
	fol.ExpiresAt = fol.CreatedAt.Add(time.Duration(fol.Expiry) * time.Hour)

	db := L.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Insert(&fol)
	if err != nil {
		panic(err)
	}
}

func FindFolders(usr *User) *[]Folder {
	db := L.NewDb("folders")
	defer db.Session.Close()
	folders := []Folder{}
	err := db.Collection.Find(bson.M{"userId": usr.Id}).All(&folders)
	if err != nil {
		panic(err)
	}
	return &folders
}

func FindFolderById(id string) (*Folder, error) {
	db := L.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	err := db.Collection.Find(bson.M{"id": id}).One(&fol)
	if err != nil {
		return nil, err
	}
	return &fol, nil
}

func (fol *Folder) Update(usr *User) {
	now := time.Now()
	m := bson.M{"updatedAt": now}

	if fol.Name != "" {
		m["name"] = fol.Name
	}

	if fol.Expiry != 0 {
		m["expiry"] = fol.Expiry
		m["expiresAt"] = now.Add(time.Duration(fol.Expiry) * time.Hour)
		fol.ExpiresAt = m["expiresAt"].(time.Time)
	}

	db := L.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Update(bson.M{"id": fol.Id, "userId": usr.Id},
		bson.M{"$set": m})
	if err != nil {
		panic(err)
	}
}

func DeleteFolder(id string, usr *User) {
	db := L.NewDb("folders")
	defer db.Session.Close()
	err := db.Collection.Remove(bson.M{"id": id, "userId": usr.Id})
	if err != nil {
		panic(err)
	}

	// Delete folder from file system
	err = os.RemoveAll(path.Join("/tmp", id)) // TODO: Read from config
	if err != nil {
		panic(err)
	}
}
