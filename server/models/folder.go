package models

import (
	"os"
	"path"
	"time"

	L "github.com/labstack/lytup/server/lytup"
	U "github.com/labstack/lytup/server/utils"
	"gopkg.in/mgo.v2/bson"
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

func (fol *Folder) Create() error {
	fol.Id = U.RandomString(7)
	fol.Files = []*File{}
	fol.CreatedAt = time.Now()
	fol.UpdatedAt = fol.CreatedAt
	fol.ExpiresAt = fol.CreatedAt.Add(time.Duration(fol.Expiry) * time.Hour)

	db := L.NewDb("folders")
	defer db.Session.Close()
	if err := db.Collection.Insert(&fol); err != nil {
		return err
	}
	return nil
}

func FindFolders(usr *User) ([]Folder, error) {
	db := L.NewDb("folders")
	defer db.Session.Close()
	folders := []Folder{}
	if err := db.Collection.Find(bson.M{"userId": usr.Id}).All(&folders); err != nil {
		return nil, err
	}
	return folders, nil
}

func FindFolderById(id string) (*Folder, error) {
	db := L.NewDb("folders")
	defer db.Session.Close()
	fol := Folder{}
	if err := db.Collection.Find(bson.M{"id": id}).One(&fol); err != nil {
		return nil, err
	}
	return &fol, nil
}

func (fol *Folder) Update(usr *User) error {
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
	if err := db.Collection.Update(
		bson.M{"id": fol.Id, "userId": usr.Id},
		bson.M{"$set": m}); err != nil {
		return err
	}
	return nil
}

func DeleteFolder(id string, usr *User) error {
	db := L.NewDb("folders")
	defer db.Session.Close()
	if err := db.Collection.Remove(bson.M{"id": id, "userId": usr.Id}); err != nil {
		return err
	}

	// Delete folder from file system
	if err := os.RemoveAll(path.Join("/tmp", id)); err != nil { // TODO: Read from config
		return err
	}

	return nil
}
