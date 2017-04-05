package gohttp

import (
    "net/http"
    "html/template"
)

func renderFile(responsewriter http.ResponseWriter, tpl string, data interface{}, code int) error {
    responsewriter.WriteHeader(code)
    t := template.Must(template.ParseFiles(tpl))
    return t.Execute(responsewriter, data)
}

func renderString(responsewriter http.ResponseWriter,name, str string, data interface{}, code int) error {
    responsewriter.WriteHeader(code)
    t := template.Must(template.New(name).Parse(str))
    return t.Execute(responsewriter, data)
}
