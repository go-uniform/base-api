package service

import (
	"net/http"
	"time"
)

// A package shorthand for map[string]interface{}
type M map[string]interface{}

// A package shorthand for map[string]string
type P map[string]string

// A model that encapsulates a bind
type Bind struct {
	Timeout         time.Duration
	Path            string
	Method          string
	Extract         func(r *http.Request) P
	ValidateRequest func(request M) M
	ConvertResponse func(response interface{}) []byte
	Permissions     []string
}