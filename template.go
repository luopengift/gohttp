package gohttp

import (
	"html/template"
	"strings"

	"github.com/luopengift/log"
)

// Template template
type Template map[string]*template.Template

// InitTemplate init template
func InitTemplate(webpath string) *Template {
	template := &Template{}
	if err := template.Lookup(webpath); err != nil {
		log.GetLogger("gohttp").Warn("init Template: %v", err)
	}
	return template
}

// Add string to template
func (t *Template) Add(name, tpl string) {
	(*t)[name] = template.Must(template.New(name).Parse(tpl))
}

// AddFile to template
func (t *Template) AddFile(tpl string) {
	tfile, err := template.ParseFiles(tpl)
	if err != nil {
		panic(err)
	}
	(*t)[tpl] = tfile
}

// Lookup walk around path
func (t *Template) Lookup(path string) error {
	files, err := WalkDir(path, func(path string) bool {
		return strings.HasSuffix(path, ".tpl")
	})
	if err != nil {
		return err
	}
	for _, file := range files {
		t.AddFile(file)
	}
	return nil
}
