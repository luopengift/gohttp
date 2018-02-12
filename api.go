package gohttp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// ApiOutput is sturct data need responsed.
type ApiOutput struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Err  error       `json:"err"`
	Data interface{} `json:"data"`
}

// MarshalJSON rewrite format to json, implement json.Marshaler interface.
func (api ApiOutput) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteByte('{')
	fmt.Fprintf(&buf, `"code":%d,`, api.Code)
	fmt.Fprintf(&buf, `"msg":"%s",`, api.Msg)
	fmt.Fprintf(&buf, `"err":"%v",`, api.Err)
	fmt.Fprintf(&buf, `"data":`)
	b, err := json.Marshal(api.Data)
	if err != nil {
		return nil, err
	}
	_, err = buf.Write(b)
	buf.WriteByte('}')
	return buf.Bytes(), err
}

func (api *ApiOutput) Set(code int, msg string) {
	api.Code = code
	api.Msg = msg
}

func (api *ApiOutput) Error() string {
	return fmt.Sprintf("%d: %s", api.Code, api.Msg)
}

func (api *ApiOutput) String() string {
	return fmt.Sprintf("%d: %s|%s", api.Code, api.Msg, api.Err.Error())
}

func (api *ApiOutput) Detail() string {
	return fmt.Sprintf("%d: %s\n%s\n%v", api.Code, api.Msg, api.Err.Error(), api.Data)
}

// ApiHandler designed for http api. It can used easy.
type ApiHandler struct {
	ApiOutput
	HttpHandler
}

func (ctx *ApiHandler) Finish() {
	ctx.Output(ctx.ApiOutput)
}

// RegistErrCode regist api error code and msg.
func RegistErrCode(code int, msg string) error {
	if value, ok := ErrMap[code]; !ok {
		return fmt.Errorf("code[%d] is used, msg=%s", code, value)
	}
	ErrMap[code] = msg
	return nil
}

// ErrMap
var ErrMap = map[int]string{
	0:    "success",
	1000: "success",
	1001: "fail",
	1002: "",
}
