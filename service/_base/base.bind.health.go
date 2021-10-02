package _base

import "net/http"

func init() {
	Bind("api.health", http.MethodGet, "/health", nil, nil, nil)
}