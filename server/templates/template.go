package templates

import (
	"fmt"
	"html/template"
	"io"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/golang/glog"
	L "github.com/labstack/lytup/server/lytup"
)

var (
	registry = map[string]*template.Template{}
	params   = map[string]string{ // Default parameters
		"hostname": L.Config.Hostname,
	}
)

func init() {
	box := rice.MustFindBox("emails")
	baseStr := box.MustString("base.html")

	err := box.Walk("", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			glog.Fatal(err)
		}

		// Skip directories and base html
		if info.IsDir() || path == "base.html" {
			return nil
		}

		baseTpl := template.Must(template.New("_").Parse(baseStr))
		emailStr := box.MustString(path)
		emailTpl := template.Must(baseTpl.Parse(emailStr))

		// Add to registry
		registry[path] = emailTpl

		return nil
	})

	if err != nil {
		glog.Fatal(err)
	}
}

func ExecuteTemplate(wr io.Writer, name string, data map[string]string) error {
	tpl, ok := registry[name]
	if !ok {
		return fmt.Errorf("template <%s> not found", name)
	}
	for k, v := range data {
		params[k] = v
	}
	return tpl.Execute(wr, params)
}

// https://github.com/drone/drone/blob/master/pkg/template/template.go
