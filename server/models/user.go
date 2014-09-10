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

func FindUserById(id string) (*User, error) {
	db := L.NewDb("users")
	defer db.Session.Close()
	usr := &User{}
	if err := db.Collection.FindId(bson.ObjectIdHex(id)).One(usr); err != nil {
		if err == mgo.ErrNotFound {
			return nil, L.ErrUserNotFound
		}
		return nil, err
	}
	return usr, nil
}

func FindUserByEmail(email string) (*User, error) {
	db := L.NewDb("users")
	defer db.Session.Close()
	usr := &User{}
	if err := db.Collection.Find(bson.M{"email": email}).One(usr); err != nil {
		if err == mgo.ErrNotFound {
			return nil, L.ErrEmailNotFound
		}
		return nil, err
	}
	return usr, nil
}

func (usr *User) Create() error {
	usr.Id = bson.NewObjectId()
	usr.Salt = U.RandomString(16)
	usr.PasswordHash = U.HashPassword(usr.Password, usr.Salt)
	usr.EmailVerified = false
	usr.CreatedAt = time.Now()
	usr.UpdatedAt = usr.CreatedAt

	db := L.NewDb("users")
	defer db.Session.Close()
	if err := db.Collection.Insert(usr); err != nil {
		if mgo.IsDup(err) {
			return L.ErrEmailIsRegistered
		}
	}

	// Create email verification key with expiry
	r := L.Redis()
	defer r.Close()
	key := U.RandomString(32)
	k := "verify:email:" + key
	r.Send("MULTI")
	r.Send("SET", k, usr.Id.Hex())
	r.Send("EXPIRE", k, L.Config.VerifyEmailExpiry)
	if _, err := r.Do("EXEC"); err != nil {
		return err
	}

	// Send 'verify email' email
	m := map[string]string{
		"name":  usr.FirstName,
		"email": usr.Email,
		"key":   key,
	}
	go U.EmailVerifyEmail(m)

	return nil
}

func (usr *User) Find() error {
	usr, err := FindUserById(usr.Id.Hex())
	return err
}

func (usr *User) Login(verify bool) error {
	pwd := usr.Password // Grab the password before it's lost

	if verify {
		db := L.NewDb("users")
		defer db.Session.Close()
		if err := db.Collection.Find(bson.M{"email": usr.Email}).One(usr); err != nil {
			if err == mgo.ErrNotFound {
				return L.ErrLogin
			}
			return err
		}

		if !bytes.Equal(usr.PasswordHash, U.HashPassword(pwd, usr.Salt)) {
			return L.ErrLogin
		}
	}

	token := jwt.New(jwt.GetSigningMethod("HS256"))
	token.Claims["iss"] = "lytup"
	token.Claims["sub"] = usr.Id
	token.Claims["exp"] = time.Now().Add(120 * time.Hour).Unix()
	var err error
	if usr.Token, err = token.SignedString([]byte(L.Config.Key)); err != nil {
		return err
	}

	return nil
}

func (usr *User) Update() error {
	now := time.Now()
	m := bson.M{"updatedAt": now}

	if usr.FirstName != "" {
		m["firstName"] = usr.FirstName
	}
	if usr.LastName != "" {
		m["lastName"] = usr.LastName
	}
	if usr.Email != "" {
		m["email"] = usr.Email
	}
	if usr.Password != "" {
		m["salt"] = U.RandomString(16)
		m["password"] = U.HashPassword(usr.Password, usr.Salt)
	}

	db := L.NewDb("users")
	defer db.Session.Close()
	if err := db.Collection.Update(
		bson.M{"id": usr.Id},
		bson.M{"$set": m}); err != nil {
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
