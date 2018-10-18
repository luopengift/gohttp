package gohttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// APIOutput is sturct data need responsed.
type APIOutput struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Err  error       `json:"err"`
	Data interface{} `json:"data"`
}

// MarshalJSON rewrite format to json, implement json.Marshaler interface.
func (api APIOutput) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	fmt.Fprintf(&buf, `"code":%d,`, api.Code)
	fmt.Fprintf(&buf, `"msg":"%s",`, api.Msg)
	fmt.Fprintf(&buf, `"err":"%v",`, strings.Replace(fmt.Sprintf("%v", api.Err), "\"", " ", -1))
	fmt.Fprintf(&buf, `"data":`)
	b, err := json.Marshal(api.Data)
	if err != nil {
		return nil, fmt.Errorf("%v, %v", err, api)
	}
	_, err = buf.Write(b)

	buf.WriteByte('}')

	bytes := buf.Bytes()
	for index, char := range bytes {
		if char == '\t' || char == '\n' || char == '\r' {
			bytes[index] = ' '
		}
	}
	return bytes, err
}

// Set set msg
func (api *APIOutput) Set(code int, msg string, errs ...error) {
	api.Code = code
	api.Msg = msg
	if len(errs) > 0 {
		api.Err = errs[0]
	}
}

func (api *APIOutput) String() string {
	return fmt.Sprintf("%d: %s|%s", api.Code, api.Msg, api.Err.Error())
}

// Detail detail
func (api *APIOutput) Detail() string {
	return fmt.Sprintf("%d: %s\n%s\n%v", api.Code, api.Msg, api.Err.Error(), api.Data)
}

// APIHandler designed for http api. It can used easily.
type APIHandler struct {
	APIOutput
	BaseHTTPHandler
}

// Finish finish
func (ctx *APIHandler) Finish() {
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	if ctx.Err != nil && ctx.Code == 0 {
		ctx.Code = 10001
	}
	ctx.Output(ctx.APIOutput)
}

// RegistErrCode regist api error code and msg.
func RegistErrCode(code int, msg string) error {
	if value, ok := ErrMap[code]; !ok {
		return fmt.Errorf("code[%d] is used, msg=%s", code, value)
	}
	ErrMap[code] = msg
	return nil
}

// ErrMap xx
var ErrMap = map[int]string{
	0:    "success",
	1000: "success",
	1001: "fail",
	1002: "",
}
