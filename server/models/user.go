package models

import (
	"bytes"
	"time"

	"github.com/dgrijalva/jwt-go"
	. "github.com/labstack/lytup/server/lytup"
	U "github.com/labstack/lytup/server/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id            bson.ObjectId `json:"id" bson:"_id"`
	FirstName     string        `json:"firstName" bson:"firstName"`
	LastName      string        `json:"lastName" bson:"lastName"`
	Email         string        `json:"email" bson:"email"`
	Password      string        `json:"password,omitempty" bson:"-"`
	PasswordHash  []byte        `json:"-" bson:"password"`
	Token         string        `json:"token,omitempty" bson:"-"`
	CreatedAt     time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time     `json:"updatedAt" bson:"updatedAt"`
	EmailVerified bool          `json:"emailVerified" bson:"emailVerified"`
	Salt          string        `json:"-" bson:"salt"`
}

func (usr *User) FindById() error {
	db := NewDb("users")
	defer db.Close()
	if err := db.Collection.FindId(usr.Id).One(usr); err != nil {
		if err == mgo.ErrNotFound {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (usr *User) FindByEmail() error {
	db := NewDb("users")
	defer db.Close()
	if err := db.Collection.Find(bson.M{"email": usr.Email}).One(usr); err != nil {
		if err == mgo.ErrNotFound {
			return ErrEmailNotFound
		}
		return err
	}
	return nil
}

func (usr *User) Create() error {
	if err := usr.validate(); err != nil {
		return err
	}

	usr.Id = bson.NewObjectId()
	usr.Salt = U.RandomString(16)
	usr.PasswordHash = U.HashPassword(usr.Password, usr.Salt)
	usr.EmailVerified = false
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = usr.CreatedAt

	db := NewDb("users")
	defer db.Close()
	if err := db.Collection.Insert(usr); err != nil {
		if mgo.IsDup(err) {
			return ErrEmailIsRegistered
		}
	}

	// Create verify email key with expiry
	code := U.RandomString(32)
	key := "verify:email:" + code
	val := usr.Id.Hex()
	if err := SetKeyWithExpiry(key, val, C.VerifyEmailExpiry); err != nil {
		return err
	}

	// Send 'verify email' email
	m := map[string]string{
		"name":  usr.FirstName,
		"email": usr.Email,
		"code":  code,
	}
	go U.EmailVerifyEmail(m)

	return nil
}

func (usr *User) Update(m bson.M) error {
	db := NewDb("users")
	defer db.Close()

	if err := db.Collection.UpdateId(
		usr.Id,
		bson.M{"$set": m}); err != nil {
		if err == mgo.ErrNotFound {
			return ErrUserNotFound
		}
		return err
	}
	return nil
}

func (usr *User) Login(verify bool) error {
	pwd := usr.Password // Grab the password before it's lost

	if verify {
		db := NewDb("users")
		defer db.Close()
		if err := db.Collection.Find(bson.M{"email": usr.Email}).One(usr); err != nil {
			if err == mgo.ErrNotFound {
				return ErrLogin
			}
			return err
		}

		if !bytes.Equal(usr.PasswordHash, U.HashPassword(pwd, usr.Salt)) {
			return ErrLogin
		}
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["iss"] = "lytup"
	token.Claims["sub"] = usr.Id
	token.Claims["exp"] = time.Now().Add(120 * time.Hour).Unix()
	var err error
	if usr.Token, err = token.SignedString([]byte(C.Key)); err != nil {
		return err
	}

	return nil
}

func (usr *User) ResetPassword() error {
	if err := usr.validatePassword(); err != nil {
		return err
	}

	salt := U.RandomString(16)
	m := bson.M{
		"salt":     salt,
		"password": U.HashPassword(usr.Password, salt),
	}
	if err := usr.Update(m); err != nil {
		return err
	}
	return nil
}

// func (usr *User) Update() error {
// 	now := time.Now()
// 	m := bson.M{"updatedAt": now}
//
// 	if usr.FirstName != "" {
// 		m["firstName"] = usr.FirstName
// 	}
// 	if usr.LastName != "" {
// 		m["lastName"] = usr.LastName
// 	}
// 	if usr.Email != "" {
// 		m["email"] = usr.Email
// 	}
// 	if usr.Password != "" {
// 		m["salt"] = U.RandomString(16)
// 		m["password"] = U.HashPassword(usr.Password, usr.Salt)
// 	}
//
// 	db := NewDb("users")
// 	defer db.Close()
// 	if err := db.Collection.Update(
// 		bson.M{"id": usr.Id},
// 		bson.M{"$set": m}); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (usr *User) Save() error {
	usr.UpdatedAt = time.Now()

	db := NewDb("users")
	defer db.Close()
	if err := db.Collection.UpdateId(usr.Id, usr); err != nil {
		return err
	}
	return nil
}

func (usr *User) validatePassword() error {
	if len(usr.Password) == 0 || len(usr.Password) <= 6 || len(usr.Password) > 32 {
		return ErrInvalidPassword
	}
	return nil
}

func (usr *User) validate() error {
	switch {
	case len(usr.FirstName) == 0, len(usr.FirstName) > 16:
		return ErrInvalidFirstName
	case len(usr.LastName) == 0, len(usr.LastName) > 16:
		return ErrInvalidLastName
	case len(usr.Email) == 0, len(usr.Email) > 32, !U.RegexpEmail.MatchString(usr.Email):
		return ErrInvalidEmail
	}
	if err := usr.validatePassword(); err != nil {
		return err
	}
	return nil
}

func (usr *User) ToRender() *User {
	usr.Password = ""
	return usr
}
