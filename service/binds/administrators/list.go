package administrators

import (
	"fmt"
	"github.com/go-uniform/uniform"
	"net/http"
	"service/service/_base"
	"strings"
)

func init() {
	_base.Bind(_base.TargetList("administrators", ""), http.MethodGet, "/administrators", func(r *http.Request) uniform.P {
		// todo: add validation rules for filters
		parameters := uniform.P{}

		parameters["-pageSize"] = r.Header.Get("Page-Size")
		if parameters["-pageSize"] == "" {
			parameters["-pageSize"] = "20"
		}

		parameters["-pageIndex"] = r.Header.Get("Page-Index")
		if parameters["-pageIndex"] == "" {
			parameters["-pageIndex"] = "1"
		}

		for key, values := range r.URL.Query() {
			if values == nil || len(values) <= 0 {
				continue
			}

			switch strings.ToLower(key) {
			default:
				panic(fmt.Sprintf("unexpected filter '%s' encountered", key))
			case "id":
				parameters["id"] = strings.Join(values, ",")
				break
			case "firstName":
				parameters["firstName"] = strings.Join(values, ",")
				break
			case "lastName":
				parameters["lastName"] = strings.Join(values, ",")
				break
			case "email":
				parameters["email"] = strings.Join(values, ",")
				break
			case "mobile":
				parameters["mobile"] = strings.Join(values, ",")
				break
			case "createdAt":
				parameters["createdAt"] = strings.Join(values, ",")
				break
			case "modifiedAt":
				parameters["modifiedAt"] = strings.Join(values, ",")
				break
			case "deletedAt":
				parameters["deletedAt"] = strings.Join(values, ",")
				break
			case "-showDeleted":
				parameters["-showDeleted"] = strings.Join(values, ",")
				break
			case "-text":
				parameters["-text"] = strings.Join(values, ",")
				break
			case "-order":
				parameters["-order"] = strings.Join(values, ",")
				break
			case "-filterToAllowable":
				parameters["-filterToAllowable"] = strings.Join(values, ",")
				break
			}
		}

		return parameters
	}, nil, nil)
}