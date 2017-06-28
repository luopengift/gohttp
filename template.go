package gohttp

import (
	"html/template"
//	"net/http"
)

type Template map[string]*template.Template

func InitTemplate() *Template {
	return &Template{}
}

// add string to template
func (t *Template) Add(name, tpl string) {
    (*t)[name] = template.Must(template.New(name).Parse(tpl))
}

// add file to template
func (t *Template) AddFile(tpl string) {
    tfile, err := template.ParseFiles(tpl)
    if err != nil {
        panic(err)
    }
    (*t)[tpl] = tfile
}
