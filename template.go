package gohttp

import (
    "net/http"
    "html/template"
)

//全局模版
var TemplateMap = map[string]*template.Template{}


func AddTemplate(name,tpl string) {
    TemplateMap[name] = template.Must(template.New(name).Parse(tpl))
}


func renderFile(responsewriter http.ResponseWriter, tpl string, data interface{}, code int) error {
    responsewriter.WriteHeader(code)
    t := template.Must(template.ParseFiles(tpl))
    return t.Execute(responsewriter, data)
}

func renderString(responsewriter http.ResponseWriter,name string, data interface{}, code int) error {
    responsewriter.WriteHeader(code)
    if t,ok := TemplateMap[name]; ok {
        return t.Execute(responsewriter, data)
    }
    return nil
}


