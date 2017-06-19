package gohttp

import (
	"encoding/json"
	"strings"
	"unsafe"
)

func BytesToJson(b []byte) (m map[string]interface{}, err error) {
    err = json.Unmarshal(b, &m)
    return
}

// Bytes2String直接转换底层指针，两者指向的相同的内存，改一个另外一个也会变。
// 效率是string([]byte{})的百倍以上，且转换量越大效率优势越明显。
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// String2Bytes直接转换底层指针，两者指向的相同的内存，改一个另外一个也会变。
// 效率是string([]byte{})的百倍以上，且转换量越大效率优势越明显。
// 转换之后若没做其他操作直接改变里面的字符，则程序会崩溃。
// 如 b:=String2bytes("xxx"); b[1]='d'; 程序将panic。
func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes(v interface{}) ([]byte, error) {
	switch v.(type) {
	case string:
		return []byte(v.(string)), nil
	case []byte:
		return v.([]byte), nil
	default:
		return json.Marshal(v)
	}
}

func String(v interface{}) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	case []byte:
		return string(v.([]byte)), nil
	default:
		b, err := json.Marshal(v)
		return string(b), err
	}

}

func hasSuffixs(s string, suffixs ...string) bool {
	for _, suffix := range suffixs {
		if ok := strings.HasSuffix(s, suffix); ok {
			return true
		}
	}
	return false
}
