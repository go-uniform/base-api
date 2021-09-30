package service

import (
	"fmt"
	"net/http"
)

// A private package level variable that contains all endpoint bindings
var bindings = map[string]Bind{}

// add a endpoint to the bindings map
func bind(topic, method, path string, extract func(r *http.Request) P, validateRequest func(request M) M, convertResponse func(response interface{}) []byte, permissions ...string) {
	if _, exists := bindings[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	bindings[topic] = Bind{
		Method:          method,
		Path:            path,
		Extract:         extract,
		ValidateRequest: validateRequest,
		ConvertResponse: convertResponse,
		Permissions:     permissions,
	}
}