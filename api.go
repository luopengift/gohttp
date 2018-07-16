package gohttp

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"strings"
// )

// // ApiOutput is sturct data need responsed.
// type ApiOutput struct {
// 	Code int         `json:"code"`
// 	Msg  string      `json:"msg"`
// 	Err  error       `json:"err"`
// 	Data interface{} `json:"data"`
// }

// // MarshalJSON rewrite format to json, implement json.Marshaler interface.
// func (api ApiOutput) MarshalJSON() ([]byte, error) {
// 	var buf bytes.Buffer
// 	buf.WriteByte('{')
// 	fmt.Fprintf(&buf, `"code":%d,`, api.Code)
// 	fmt.Fprintf(&buf, `"msg":"%s",`, api.Msg)
// 	fmt.Fprintf(&buf, `"err":"%v",`, strings.Replace(fmt.Sprintf("%v", api.Err), "\"", " ", -1))
// 	fmt.Fprintf(&buf, `"data":`)
// 	b, err := json.Marshal(api.Data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_, err = buf.Write(b)

// 	buf.WriteByte('}')

// 	bytes := buf.Bytes()
// 	for index, char := range bytes {
// 		if char == '\t' || char == '\n' || char == '\r' {
// 			bytes[index] = ' '
// 		}
// 	}
// 	return bytes, err
// }

// func (api *ApiOutput) Set(code int, msg string, errs ...error) {
// 	api.Code = code
// 	api.Msg = msg
// 	if len(errs) > 0 {
// 		api.Err = errs[0]
// 	}
// }

// func (api *ApiOutput) String() string {
// 	return fmt.Sprintf("%d: %s|%s", api.Code, api.Msg, api.Err.Error())
// }

// func (api *ApiOutput) Detail() string {
// 	return fmt.Sprintf("%d: %s\n%s\n%v", api.Code, api.Msg, api.Err.Error(), api.Data)
// }

// // ApiHandler designed for http api. It can used easily.
// type ApiHandler struct {
// 	ApiOutput
// 	HttpHandler
// }

// func (ctx *ApiHandler) Finish() {
// 	ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
// 	ctx.Output(ctx.ApiOutput)
// }

// // RegistErrCode regist api error code and msg.
// func RegistErrCode(code int, msg string) error {
// 	if value, ok := ErrMap[code]; !ok {
// 		return fmt.Errorf("code[%d] is used, msg=%s", code, value)
// 	}
// 	ErrMap[code] = msg
// 	return nil
// }

// // ErrMap
// var ErrMap = map[int]string{
// 	0:    "success",
// 	1000: "success",
// 	1001: "fail",
// 	1002: "",
// }
