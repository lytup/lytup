package models

import (
	"bytes"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	L "github.com/labstack/lytup/server/lytup"
	U "github.com/labstack/lytup/server/utils"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id               bson.ObjectId `json:"id" bson:"_id"`
	Name             string        `json:"name" bson:"name"`
	Email            string        `json:"email" bson:"email"`
	Password         string        `json:"password,omitempty" bson:"-"`
	PasswordHash     []byte        `json:"-" bson:"password"`
	Token            string        `json:"token,omitempty" bson:"-"`
	CreatedAt        time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time     `json:"updatedAt" bson:"updatedAt"`
	ConfirmationCode string        `json:"-" bson:"confirmationCode"`
	Confirmed        bool          `json:"confirmed" bson:"confirmed"`
	Salt             string        `json:"-" bson:"salt"`
}

func FindUserById(id string) (*User, *L.Error) {
	db := L.NewDb("users")
	defer db.Session.Close()
	usr := &User{}
	if err := db.Collection.FindId(id).One(usr); err != nil {
		if err.Error() == "not found" {
			return nil, L.NewError("invalid username or password", L.MongoDbNotFoundError)
		}
		return nil, L.NewError(err.Error(), L.MongoDbError)
	}
	return usr, nil
}

func FindUserByConfirmationCode(code string) (*User, *L.Error) {
	db := L.NewDb("users")
	defer db.Session.Close()
	usr := &User{}
	if err := db.Collection.Find(bson.M{"confirmationCode": code}).One(usr); err != nil {
		if err.Error() == "not found" {
			return nil, L.NewError("invalid confirmation code", L.MongoDbNotFoundError)
		}
		return nil, L.NewError(err.Error(), L.MongoDbError)
	}
	return usr, nil
}

func (usr *User) Create() error {
	usr.Id = bson.NewObjectId()
	usr.Salt = U.RandomString(16)
	usr.PasswordHash = U.HashPassword(usr.Password, usr.Salt)
	usr.ConfirmationCode = U.RandomString(32)
	usr.Confirmed = false
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = usr.CreatedAt

	db := L.NewDb("users")
	defer db.Session.Close()
	if err := db.Collection.Insert(usr); err != nil {
		return err
	}
	return nil
}

func (usr *User) Login() *L.Error {
	pwd := usr.Password // Grab the password before it's lost

	db := L.NewDb("users")
	defer db.Session.Close()
	if err := db.Collection.Find(bson.M{"email": usr.Email}).One(usr); err != nil {
		if err.Error() == "not found" {
			return L.NewError(fmt.Sprintf("email <%s> not found", usr.Email), L.LoginError)
		}
		return L.NewError(err.Error(), L.MongoDbError)
	}

	if !bytes.Equal(usr.PasswordHash, U.HashPassword(pwd, usr.Salt)) {
		return L.NewError("invalid password", L.LoginError)
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["exp"] = time.Now().Add(120 * time.Hour).Unix()
	token.Claims["usr-id"] = usr.Id
	var err error
	if usr.Token, err = token.SignedString([]byte(L.Config.Key)); err != nil {
		return L.NewError(err.Error(), L.JwtError)
	}

	return nil
}

func (usr *User) Save() error {
	usr.UpdatedAt = time.Now()

	db := L.NewDb("users")
	defer db.Session.Close()
	if err := db.Collection.UpdateId(usr.Id, usr); err != nil {
		return err
	}
	return nil
}

func (usr *User) Render() *User {
	usr.Password = ""
	return usr
}
