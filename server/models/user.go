package models

import (
	"github.com/labstack/lytup/server/db"
	"github.com/labstack/lytup/server/utils"
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	Id             bson.ObjectId `json:"id" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Email          string        `json:"email" bson:"email"`
	Password       string        `json:"password,omitempty" bson:"-"`
	HashedPassword []byte        `json:"-" bson:"password"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	EmailVerified  bool          `json:"emailVerified" bson:"emailVerified"`
}

func (usr *User) Create() {
	usr.Id = bson.NewObjectId()
	usr.CreatedAt = time.Now()

	db := db.NewDb("users")
	defer db.Session.Close()
	err := db.Collection.Insert(usr)
	if err != nil {
		panic(err)
	}
}

func (usr *User) Login() error {
	db := db.NewDb("users")
	defer db.Session.Close()
	return db.Collection.Find(bson.M{"email": usr.Email,
		"password": utils.HashPassword([]byte(usr.Password))}).
		One(usr)
}

func (self *User) Render() *User {
	self.Password = ""
	return self
}
