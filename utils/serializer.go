package utils

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"reflect"
	"strconv"
)

type BlobEncode interface {
	Encode(val interface{}) ([]byte, error)
	Decode(data []byte, ptr interface{}) error
}

type JsonEncode struct{}

func (j JsonEncode) Encode(val interface{}) ([]byte, error) {
	if v, ok := val.([]byte); ok {
		return v, nil
	}

	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	return json.Marshal(val)
}

func (j JsonEncode) Decode(data []byte, ptr interface{}) error {
	if v, ok := ptr.(*[]byte); ok {
		*v = data
		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var x int64
			x, err := strconv.ParseInt(string(data), 10, 64)
			if err != nil {
				return err
			}
			p.SetInt(x)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var x uint64
			x, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return err
			}
			p.SetUint(x)
			return nil
		}
	}

	return json.Unmarshal(data, ptr)
}

type GobEncode struct{}

func (g GobEncode) Encode(val interface{}) ([]byte, error) {
	if v, ok := val.([]byte); ok {
		return v, nil
	}

	switch v := reflect.ValueOf(val); v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return []byte(strconv.FormatInt(v.Int(), 10)), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return []byte(strconv.FormatUint(v.Uint(), 10)), nil
	}

	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	if err := encoder.Encode(val); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (g GobEncode) Decode(data []byte, ptr interface{}) error {
	if v, ok := ptr.(*[]byte); ok {
		*v = data
		return nil
	}

	if v := reflect.ValueOf(ptr); v.Kind() == reflect.Ptr {
		switch p := v.Elem(); p.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var x int64
			x, err := strconv.ParseInt(string(data), 10, 64)
			if err != nil {
				return err
			}
			p.SetInt(x)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			var x uint64
			x, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return err
			}
			p.SetUint(x)
			return nil
		}
	}

	b := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(b)
	return decoder.Decode(ptr)
}

func Serializer(encoding BlobEncode, val interface{}) ([]byte, error) {
	return encoding.Encode(val)
}

func Deserializer(encoding BlobEncode, data []byte, ptr interface{}) error {
	return encoding.Decode(data, ptr)
}
