package _base

import (
	"github.com/go-uniform/uniform"
	"net/http"
	"time"
)

// A model that encapsulates a bind
type bind struct {
	Timeout         time.Duration
	Path            string
	Method          string
	Extract         func(r *http.Request) uniform.P
	ValidateRequest func(request uniform.M) uniform.M
	ConvertResponse func(response interface{}) []byte
	Permissions     []string
}