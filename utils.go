package gohttp

import (
	"encoding/json"
)

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
