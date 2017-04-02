package gohttp

import (
    "net/http"
    "html/template"
)

func render(responsewriter http.ResponseWriter, tpl string, data interface{}, code int) error {
    responsewriter.WriteHeader(code)
    t, err := template.ParseFiles(tpl)
    if err != nil {
        return err
    }
    return t.Execute(responsewriter, data)
}
