package gohttp

import (
	"html/template"
	//	"net/http"
)

// Template template
type Template map[string]*template.Template

// InitTemplate init template
func InitTemplate() *Template {
	return &Template{}
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
