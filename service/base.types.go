package service

import "net/http"

// A package shorthand for map[string]interface{}
type M map[string]interface{}

// A model that encapsulates a bind
type Bind struct {
	App string
	Path string
	Method string
	Topic string
	Tags map[string][]string
	Extract func(r *http.Request) (map[string]string, error)
	Convert func(data []byte) (interface{}, error)
	ConvertBack func(data []byte) ([]byte, error)
	Permissions []string
}