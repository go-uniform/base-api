package administrators

import (
	"net/http"
	"service/service/_base"
)

func init() {
	_base.Bind(_base.TargetItem("administrators", ""), http.MethodGet, "/administrators/{id}", _base.ExtractIdPathParameter, nil, nil)
}