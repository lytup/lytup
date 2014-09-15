package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	. "github.com/labstack/lytup/server/lytup"
	"github.com/labstack/lytup/server/models"
	R "github.com/labstack/lytup/server/router"
	U "github.com/labstack/lytup/server/utils"
	"gopkg.in/mgo.v2/bson"

	. "github.com/smartystreets/goconvey/convey"
)

func newUser() *models.User {
	return &models.User{
		Id:        bson.NewObjectId(),
		FirstName: "Vishal",
		LastName:  "Rana",
		Email:     "lytup.test@labstack.com",
		Password:  "password",
	}
}

func seedUser(usr *models.User, db *Db) *models.User {
	if usr == nil {
		usr = newUser()
	}
	db.Collection.Insert(usr)
	return usr
}

func removeUser(db *Db, usr *models.User) {
	db.Collection.Remove(bson.M{"email": usr.Email})
}

func TestCreateUser(t *testing.T) {
	db := NewDb("users")
	defer db.Close()

	Convey("Create User", t, func() {
		req := &http.Request{}
		rr := httptest.NewRecorder()
		ctx := R.NewContext(rr, req, nil, nil)

		Convey("It successfully creates a user", func() {
			usr := newUser()
			usrStr, _ := json.Marshal(usr)
			req.Body = ioutil.NopCloser(bytes.NewReader(usrStr))
			CreateUser(ctx)
			removeUser(db, usr)
			So(rr.Code, ShouldEqual, http.StatusCreated)
		})

		Convey("It raises an error for invalid JSON", func() {
			usrStr := `{"invalid": true}`
			req.Body = ioutil.NopCloser(strings.NewReader(usrStr))
			CreateUser(ctx)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
		})

		// Invalid first name
		Convey("It raises an error if first name is invalid", func() {
			usr := newUser()
			usr.FirstName = "Vishal Vishal Vishal"
			usrStr, _ := json.Marshal(usr)
			req.Body = ioutil.NopCloser(bytes.NewReader(usrStr))
			CreateUser(ctx)
			removeUser(db, usr)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
			So(rr.Body.String(), ShouldContainSubstring, M["invalidFirstName"])
		})

		// Invalid last name
		Convey("It raises an error if last name is invalid", func() {
			usr := newUser()
			usr.LastName = "Rana Rana Rana Rana"
			usrStr, _ := json.Marshal(usr)
			req.Body = ioutil.NopCloser(bytes.NewReader(usrStr))
			CreateUser(ctx)
			removeUser(db, usr)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
			So(rr.Body.String(), ShouldContainSubstring, M["invalidLastName"])
		})

		// Invalid email
		Convey("It raises an error if email is invalid", func() {
			usr := newUser()
			usr.Email = "bad-email"
			usrStr, _ := json.Marshal(usr)
			req.Body = ioutil.NopCloser(bytes.NewReader(usrStr))
			CreateUser(ctx)
			removeUser(db, usr)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
			So(rr.Body.String(), ShouldContainSubstring, M["invalidEmail"])
		})

		// Email is already registered
		Convey("It raises an error if email is already registered", func() {
			usr := seedUser(nil, db)
			usrStr, _ := json.Marshal(usr)
			req.Body = ioutil.NopCloser(bytes.NewReader(usrStr))
			CreateUser(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusConflict)
			So(rr.Body.String(), ShouldContainSubstring, M["emailIsRegisteredError"])
		})
	})
}

func TestVerifyEmail(t *testing.T) {
	db := NewDb("users")
	defer db.Close()

	Convey("Verify Email", t, func() {
		req := &http.Request{}
		rr := httptest.NewRecorder()
		code := U.RandomString(32)
		p := httprouter.Params{
			httprouter.Param{"code", code},
		}
		ctx := R.NewContext(rr, req, p, nil)

		Convey("It successfully verifies an email", func() {
			usr := seedUser(nil, db)
			key := "verify:email:" + code
			val := usr.Id.Hex()
			SetKeyWithExpiry(key, val, C.VerifyEmailExpiry)
			VerifyEmail(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, 200)
			So(rr.Body.String(), ShouldContainSubstring, M["verifyEmailSuccess"])
		})

		Convey("It raises an error for bad email verification key", func() {
			usr := seedUser(nil, db)
			key := "verify:email:" + code
			val := usr.Id.Hex()
			SetKeyWithExpiry(key, val, C.VerifyEmailExpiry)
			p = httprouter.Params{
				httprouter.Param{"code", "bad-code"},
			}
			ctx = R.NewContext(rr, req, p, nil)
			VerifyEmail(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
			So(rr.Body.String(), ShouldContainSubstring, M["verifyEmailError"])
		})

		Convey("It raises an error if email is already verified", func() {
			usr := newUser()
			usr.EmailVerified = true
			seedUser(usr, db)
			key := "verify:email:" + code
			val := usr.Id.Hex()
			SetKeyWithExpiry(key, val, C.VerifyEmailExpiry)
			VerifyEmail(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
			So(rr.Body.String(), ShouldContainSubstring, M["verifyEmailError"])
		})
	})
}

func TestForgotPassword(t *testing.T) {
	db := NewDb("users")
	defer db.Close()

	Convey("Forgot Password", t, func() {
		req := &http.Request{}
		rr := httptest.NewRecorder()
		ctx := R.NewContext(rr, req, nil, nil)

		Convey("It successfully accepts forgot password request", func() {
			usr := seedUser(nil, db)
			usrStr := `{"email": "lytup.test@labstack.com"}`
			req.Body = ioutil.NopCloser(strings.NewReader(usrStr))
			ForgotPassword(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusOK)
		})

		Convey("It raises an error if email is not found", func() {
			usr := seedUser(nil, db)
			usrStr := `{"email": "bad-email"}`
			req.Body = ioutil.NopCloser(strings.NewReader(usrStr))
			ForgotPassword(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusNotFound)
			So(rr.Body.String(), ShouldContainSubstring, M["emailNotFoundError"])
		})
	})
}

func TestResetPassword(t *testing.T) {
	db := NewDb("users")
	defer db.Close()

	Convey("Reset Password", t, func() {
		req := &http.Request{}
		rr := httptest.NewRecorder()
		code := U.RandomString(32)
		p := httprouter.Params{
			httprouter.Param{"code", code},
		}
		ctx := R.NewContext(rr, req, p, nil)

		Convey("It successfully resets the password", func() {
			usr := seedUser(nil, db)
			key := "reset:pass:" + code
			val := usr.Id.Hex()
			SetKeyWithExpiry(key, val, C.PasswordResetExpiry)
			usrStr := `{"password": "password2"}`
			req.Body = ioutil.NopCloser(strings.NewReader(usrStr))
			ResetPassword(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusOK)
		})

		Convey("It raises an error if password is invalid", func() {
			usr := seedUser(nil, db)
			key := "reset:pass:" + code
			val := usr.Id.Hex()
			SetKeyWithExpiry(key, val, C.PasswordResetExpiry)
			usrStr := `{"password": "pwd"}`
			req.Body = ioutil.NopCloser(strings.NewReader(usrStr))
			ResetPassword(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
		})

		Convey("It raises an error for bad password reset key", func() {
			usr := seedUser(nil, db)
			key := "reset:pass:" + code
			val := usr.Id.Hex()
			SetKeyWithExpiry(key, val, C.PasswordResetExpiry)
			p = httprouter.Params{
				httprouter.Param{"code", "bad-code"},
			}
			ctx = R.NewContext(rr, req, p, nil)
			usrStr := `{"password": "password2"}`
			req.Body = ioutil.NopCloser(strings.NewReader(usrStr))
			ResetPassword(ctx)
			db.Collection.RemoveId(usr.Id)
			So(rr.Code, ShouldEqual, http.StatusBadRequest)
		})
	})
}

// func TestFindUser(t *testing.T) {
// 	db := NewDb("users")
// 	defer db.Close()
//
// 	Convey("Find User", t, func() {
// 		req := &http.Request{}
// 		rr := httptest.NewRecorder()
// 		ctx := R.NewContext(rr, req, nil, nil)
//
// 		Convey("It successfully finds the user", func() {
// 			usr := seedUser(nil, db)
// 			ctx.User = usr
// 			FindUser(ctx)
// 			db.Collection.RemoveId(usr.Id)
// 			So(rr.Code, ShouldEqual, http.StatusOK)
// 		})
//
// 		Convey("It raises an error if the user is not found", func() {
// 			usr := seedUser(nil, db)
// 			ctx.User = newUser()
// 			FindUser(ctx)
// 			db.Collection.RemoveId(usr.Id)
// 			So(rr.Code, ShouldEqual, http.StatusNotFound)
// 		})
// 	})
// }
