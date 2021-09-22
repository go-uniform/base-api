package service

import (
	"net/http"
	"time"
)

// A package shorthand for map[string]interface{}
type M map[string]interface{}

// A model that encapsulates a bind
type Bind struct {
	Timeout time.Duration
	Path string
	Method string
	Extract func(r *http.Request) (map[string]string, error)
	ConvertRequest func(data []byte) (interface{}, error)
	ConvertResponse func(data []byte) ([]byte, error)
	Permissions []string
}