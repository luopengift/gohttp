package gohttp

import (
	"html/template"
	"net/http"
)

//全局模版
var TemplateMap = map[string]*template.Template{}

func AddTemplate(name, tpl string) {
	TemplateMap[name] = template.Must(template.New(name).Parse(tpl))
}

func renderFile(responsewriter http.ResponseWriter, tpl string, data interface{}) (int, error) {
	t, err := template.ParseFiles(tpl)
	if err != nil {
		http.Error(responsewriter, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return http.StatusNotFound, err
	}
	return http.StatusOK, t.Execute(responsewriter, data)
}

func renderString(responsewriter http.ResponseWriter, name string, data interface{}) error {
	if t, ok := TemplateMap[name]; ok {
		return t.Execute(responsewriter, data)
	}
	return nil
}

type Template map[string]*template.Template

func InitTemplate() *Template {
	return new(Template)
}

func (t *Template) Add(name, tpl string) {
	(*t)[name] = template.Must(template.New(name).Parse(tpl))
}

func (t *Template) renderFile(tpl string, data interface{}) {}
