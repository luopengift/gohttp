package gohttp

import (
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/luopengift/types"
)

// Context gohttp context
type Context struct {
	// Context Std lib
	//context.Context
	*response
	*http.Request

	*Application
	match map[string]string
	body  []byte
}

func (ctx *Context) reset() {
	//ctx.Context = nil
	ctx.response = nil
	ctx.Request = nil
}

func (ctx *Context) init(resp http.ResponseWriter, req *http.Request) {
	ctx.response = newResponseWriter(resp)
	ctx.Request = req
}

// RemoteAddr remote addr for request, copied from net/url.stripPort
func (ctx *Context) RemoteAddr() string {
	colon := strings.IndexByte(ctx.Request.RemoteAddr, ':')
	if colon == -1 {
		return ctx.Request.RemoteAddr
	}
	if i := strings.IndexByte(ctx.Request.RemoteAddr, ']'); i != -1 {
		return strings.TrimPrefix(ctx.Request.RemoteAddr[:i], "[")
	}
	return ctx.Request.RemoteAddr[:colon]
}

// GetCookies get cookies
func (ctx *Context) GetCookies() []*http.Cookie {
	return ctx.Request.Cookies()
}

// GetCookie get cookie
func (ctx *Context) GetCookie(name string) string {
	cookie, err := ctx.Request.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

// HTTPError response Http Error
func (ctx *Context) HTTPError(msg string, code int) {
	ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
	ctx.ResponseWriter.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.output([]byte(msg), code)
}

// Redirect response redirect
func (ctx *Context) Redirect(url string, code int) {
	http.Redirect(ctx.ResponseWriter, ctx.Request, url, code)
}

func (ctx *Context) output(response []byte, code int) {
	if ctx.Finished() {
		ctx.Log.Warnf("[warn] output called twice!")
		return
	}
	ctx.WriteHeader(code)
	ctx.ResponseWriter.Write(response)
}

// Output xx
func (ctx *Context) Output(v interface{}, code ...int) {
	response, err := types.ToBytes(v)
	if err != nil {
		panic(err)
	}
	if len(code) == 0 {
		ctx.output(response, http.StatusOK)
	} else {
		ctx.output(response, code[0])
	}
}

// parseArguments parse and handler arguments
func (ctx *Context) parseArgs() error {
	// parse form automatically
	if err := ctx.Request.ParseForm(); err != nil {
		return err
	}
	//ctx.match = match
	return nil
}

// GetMatch fetch match argument named by <name>, null is default value defined by user
func (ctx *Context) GetMatch(name string, null string) string {
	value, ok := ctx.match[name]
	if !ok {
		return null
	}
	return value
}

// GetMatchArgs get match args
func (ctx *Context) GetMatchArgs() map[string]string {
	return ctx.match
}

// GetQuery fetch query argument named by <name>, null is default value defined by user
func (ctx *Context) GetQuery(name string, null string) string {
	value, ok := ctx.Request.Form[name]
	if !ok {
		return null
	}
	return value[0]
}

// GetQueryArgs get query args
func (ctx *Context) GetQueryArgs() map[string][]string {
	return ctx.Request.Form
}

// GetBody fetch body argument named by <name>
func (ctx *Context) GetBody(name string) interface{} {
	ctx.GetBodyArgs()
	if body, err := types.BytesToMap(ctx.body); err != nil {
		panic(err)
	} else {
		return body[name]
	}
}

// GetBodyArgs fetch body arguments
func (ctx *Context) GetBodyArgs() []byte {
	body, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		panic(err)
	}
	if len(body) != 0 {
		ctx.body = body
	}
	return ctx.body
}

// Download file download response by file path.
func (ctx *Context) Download(file string) {
	if ctx.Finished() {
		ctx.Log.Errorf("HttpHandler is end!")
		return
	}
	f, err := os.Stat(file)
	if err != nil {
		ctx.HTTPError(http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	if f.IsDir() {
		ctx.HTTPError(http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/octet-stream")
	ctx.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename="+path.Base(file))
	ctx.ResponseWriter.Header().Set("Content-Transfer-Encoding", "binary")
	ctx.ResponseWriter.Header().Set("Expires", "0")
	ctx.ResponseWriter.Header().Set("Cache-Control", "must-revalidate")
	ctx.ResponseWriter.Header().Set("Pragma", "public")
	http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
	return
}

// GetForm formdata, Content-Type must be multipart/form-data.
// TODO: RemoveAll removes any temporary files associated with a Form.
func (ctx *Context) GetForm() (map[string]string, map[string]*multipart.FileHeader, error) {
	reader, err := ctx.Request.MultipartReader()
	if err != nil {
		return nil, nil, err
	}
	form, err := reader.ReadForm(10000)
	if err != nil {
		return nil, nil, err
	}
	values := make(map[string]string)
	for k, v := range form.Value {
		if len(v) > 0 {
			values[k] = v[0]
		}
	}
	files := make(map[string]*multipart.FileHeader)
	for k, v := range form.File {
		if len(v) > 0 {
			files[k] = v[0]
		}
	}
	return values, files, nil
}

//SaveFile save file to disk
func (ctx *Context) SaveFile(fh *multipart.FileHeader, path string, name ...string) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var filename string
	if len(name) == 0 {
		filename = fh.Filename
	} else {
		filename = name[0]
	}
	filepath := filepath.Join(path, filename)
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return filepath, err
}

// RecvFile recv file
func (ctx *Context) RecvFile(name string, path string) (string, error) {
	file, head, err := ctx.Request.FormFile(name)
	if err != nil {
		return "", err
	}
	defer file.Close()
	filepath := filepath.Join(path, head.Filename)
	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	return filepath, err
}

// Render render template no cache
func (ctx *Context) Render(tpl string, data interface{}) {
	path := filepath.Join(ctx.Config.WebPath, tpl)
	t, err := template.ParseFiles(path)
	if err != nil {
		ctx.HTTPError(toHTTPError(err))
		return
	}
	ctx.render(t, data)
}

// render html data to client
func (ctx *Context) render(tpl *template.Template, data interface{}) {
	ctx.WriteHeader(http.StatusOK)
	tpl.Execute(ctx.ResponseWriter, data)
}

// Exec inplement HandlerHTTP interface!
// func (ctx *Context) Exec(context *Context) {
// 	ctx.Info("Exec func implement HandlerHTTP intertface!!!")
// }

// JSON response json data
func (ctx *Context) JSON(v interface{}, code ...int) {

}

// XML response xml data
func (ctx *Context) XML(v interface{}, code ...int) {

}

// HTML response html data
func (ctx *Context) HTML(v interface{}, code ...int) {

}

// HEAD method
// func (ctx *Context) HEAD() {}

// GET method, must be overwrite
// func (ctx *Context) GET() {
// 	ctx.HTTPError(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// }

// POST method, must be overwrite
// func (ctx *Context) POST() {
// 	ctx.HTTPError(http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
// }
