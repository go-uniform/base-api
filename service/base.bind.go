package service

import (
	"fmt"
	"net/http"
)

// A private package level variable that contains all endpoint bindings
var bindings = map[string]Bind{}

// add a endpoint to the bindings map
func bind(topic, method, path string, extract func(r *http.Request) (map[string]string, error), convertRequest func(data []byte) (interface{}, error), convertResponse func(data []byte) ([]byte, error), permissions ...string) {
	if _, exists := bindings[topic]; exists {
		panic(fmt.Sprintf("topic '%s' has already been subscribed", topic))
	}
	bindings[topic] = Bind{
		Method: method,
		Path: path,
		Extract: extract,
		ConvertRequest: convertRequest,
		ConvertResponse: convertResponse,
		Permissions: permissions,
	}
}