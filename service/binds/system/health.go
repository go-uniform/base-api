package system

import (
	"net/http"
	"service/service/_base"
)

func init() {
	_base.Bind(_base.TargetSystem("health"), http.MethodGet, "/health", nil, nil, nil)
}