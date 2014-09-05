package models

import (
	"bytes"
	"time"

	"github.com/dgrijalva/jwt-go"
	L "github.com/labstack/lytup/server/lytup"
	U "github.com/labstack/lytup/server/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id           bson.ObjectId `json:"id" bson:"_id"`
	FirstName    string        `json:"firstName" bson:"firstName"`
	LastName     string        `json:"lastName" bson:"lastName"`
	Email        string        `json:"email" bson:"email"`
	Password     string        `json:"password,omitempty" bson:"-"`
	PasswordHash []byte        `json:"-" bson:"password"`
	Token        string        `json:"token,omitempty" bson:"-"`
	CreatedAt    time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt" bson:"updatedAt"`
	Confirmed    bool          `json:"confirmed" bson:"confirmed"`
	Salt         string        `json:"-" bson:"salt"`
}

func FindUserById(id string) (*User, error) {
	db := L.NewDb("users")
	defer db.Session.Close()
	usr := &User{}
	if err := db.Collection.FindId(bson.ObjectIdHex(id)).One(usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func FindUserByConfirmationCode(code string) (*User, error) {
	db := L.NewDb("users")
	defer db.Session.Close()
	usr := &User{}
	if err := db.Collection.Find(bson.M{"confirmationCode": code}).One(usr); err != nil {
		return nil, err
	}
	return usr, nil
}

func (usr *User) Create() error {
	usr.Id = bson.NewObjectId()
	usr.Salt = U.RandomString(16)
	usr.PasswordHash = U.HashPassword(usr.Password, usr.Salt)
	usr.Confirmed = false
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = usr.CreatedAt

	db := L.NewDb("users")
	defer db.Session.Close()
	if err := db.Collection.Insert(usr); err != nil {
		return err
	}

	// Create user confirmation key with expiry
	r := L.Redis()
	defer r.Close()
	key := U.RandomString(32)
	r.Send("MULTI")
	r.Send("SET", key, usr.Id.Hex())
	r.Send("EXPIRE", key, L.Config.ConfirmationExpiry)
	if _, err := r.Do("EXEC"); err != nil {
		return err
	}

	// Send confirmation email
	m := map[string]string{
		"hostname": L.Config.Hostname,
		"name":     usr.FirstName,
		"email":    usr.Email,
		"key":      key,
	}
	go U.EmailConfirmation(m)

	return nil
}

func (usr *User) Find() error {
	usr, err := FindUserById(usr.Id.Hex())
	return err
}

func (usr *User) Login() error {
	pwd := usr.Password // Grab the password before it's lost

	db := L.NewDb("users")
	defer db.Session.Close()
	if err := db.Collection.Find(bson.M{"email": usr.Email}).One(usr); err != nil {
		return err
	}

	if !bytes.Equal(usr.PasswordHash, U.HashPassword(pwd, usr.Salt)) {
		return mgo.ErrNotFound
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["exp"] = time.Now().Add(120 * time.Hour).Unix()
	token.Claims["usr-id"] = usr.Id
	var err error
	if usr.Token, err = token.SignedString([]byte(L.Config.Key)); err != nil {
		return err
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
