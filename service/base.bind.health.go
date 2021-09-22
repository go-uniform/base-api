package service

import "net/http"

func init() {
	bind("api.health", http.MethodGet, "/health", nil, nil, nil)
}