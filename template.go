package gohttp

import (
	"html/template"
)

// Template template
type Template map[string]*template.Template

// InitTemplate init template
func InitTemplate(webpath string) (*Template, error) {
	template := &Template{}
	err := template.Lookup(webpath)
	return template, err
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
		return hasSuffixs(path, ".html", ".tpl")
	})
	if err != nil {
		return err
	}
	for _, file := range files {
		t.AddFile(file)
	}
	return nil
}
