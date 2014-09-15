package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	. "github.com/labstack/lytup/server/lytup"
	R "github.com/labstack/lytup/server/router"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	folJson = `{"name": "Pics", "expiry": 3600}`
)

// func newFolder() *models.Folder {
// 	return &models.Folder{
// 		Name:   "Pictures",
// 		Expiry: 3600,
// 	}
// }

func TestCreateFolder(t *testing.T) {
	db := NewDb("users")
	defer db.Close()

	Convey("Create Folder", t, func() {
		req := &http.Request{}
		rr := httptest.NewRecorder()
		ctx := R.NewContext(rr, req, nil, nil)

		Convey("It successfully creates a folder", func() {
			usr := seedUser(nil, db)
			ctx.User = usr
			req.Body = ioutil.NopCloser(strings.NewReader(folJson))
			CreateFolder(ctx)
			removeUser(db, usr)
			So(rr.Code, ShouldEqual, http.StatusCreated)
		})
	})
}
