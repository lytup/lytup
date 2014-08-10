package models

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/lytup/server/db"
	U "github.com/labstack/lytup/server/utils"
	"labix.org/v2/mgo/bson"
	"time"
)

type User struct {
	Id             bson.ObjectId `json:"id" bson:"_id"`
	Name           string        `json:"name" bson:"name"`
	Email          string        `json:"email" bson:"email"`
	Password       string        `json:"password,omitempty" bson:"-"`
	HashedPassword []byte        `json:"-" bson:"password"`
	Token          string        `json:"token,omitempty" bson:"-"`
	CreatedAt      time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt" bson:"updatedAt"`
	EmailVerified  bool          `json:"emailVerified" bson:"emailVerified"`
}

func (usr *User) Create() error {
	usr.Id = bson.NewObjectId()
	usr.CreatedAt = time.Now()

	db := db.NewDb("users")
	defer db.Session.Close()
	err := db.Collection.Insert(usr)
	if err != nil {
		return err
	}

	err = usr.Login()
	if err != nil {
		return err
	}

	return nil
}

func (usr *User) Find() error {
	db := db.NewDb("users")
	defer db.Session.Close()
	err := db.Collection.FindId(usr.Id).One(usr)
	if err != nil {
		return err
	}

	return nil
}

func (usr *User) Login() error {
	db := db.NewDb("users")
	defer db.Session.Close()
	err := db.Collection.Find(bson.M{"email": usr.Email,
		"password": U.HashPassword([]byte(usr.Password))}).
		One(usr)
	if err != nil {
		return err
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["exp"] = time.Now().Add(120 * time.Hour).Unix()
	token.Claims["usr-id"] = usr.Id
	usr.Token, err = token.SignedString(U.KEY)
	if err != nil {
		return err
	}
	return nil
}

func (usr *User) Render() *User {
	usr.Password = ""
	return usr
}
