package templates

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"

	"github.com/GeertJohan/go.rice"
)

var registry = map[string]*template.Template{}

func init() {
	box := rice.MustFindBox("emails")
	baseStr := box.MustString("base.html")

	// err := box.Walk("", func(path string, info os.FileInfo, err error) error {
	// 	if err != nil {
	// 		glog.Fatal(err)
	// 	}

	// Skip directories and base html
	// if info.IsDir() || path == "base.html" {
	// 	return nil
	// }

	path := "confirmation.html"
	baseTpl := template.Must(template.New("_").Parse(baseStr))
	emailStr := box.MustString(path)
	emailTpl := template.Must(baseTpl.Parse(emailStr))

	var b bytes.Buffer
	m := map[string]string{
		"hostname": "L.Config.Hostname",
		"name":     "usr.Name",
		"email":    "usr.Email",
		"code":     "usr.ConfirmationCode",
	}

	// Add to registry
	registry[path] = emailTpl

	ExecuteTemplate(&b, path, m)

	// 	return nil
	// })
	//
	// if err != nil {
	// 	glog.Fatal(err)
	// }
}

func ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	tpl, ok := registry[name]
	if !ok {
		return errors.New(fmt.Sprintf("template <%s> not found", name))
	}
	return tpl.Execute(wr, data)
}

// https://github.com/drone/drone/blob/master/pkg/template/template.go
