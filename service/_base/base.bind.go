package _base

import (
	"fmt"
	"github.com/go-uniform/uniform"
	"net/http"
)

// A private package level variable that contains all endpoint bindings
var Bindings = map[string]bind{}

// add a endpoint to the bindings map
func Bind(topic, method, path string, extract func(r *http.Request) uniform.P, validateRequest func(request uniform.M) uniform.M, convertResponse func(response interface{}) []byte, permissions ...string) {
	if _, exists := Bindings[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	Bindings[topic] = bind{
		Method:          method,
		Path:            path,
		Extract:         extract,
		ValidateRequest: validateRequest,
		ConvertResponse: convertResponse,
		Permissions:     permissions,
	}
}