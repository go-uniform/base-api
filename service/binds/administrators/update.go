package administrators

import (
	"github.com/go-uniform/uniform"
	"net/http"
	"service/service/_base"
)

func init() {
	_base.Bind(_base.TargetItem("administrators", "update"), http.MethodPatch, "/administrators/{id}", _base.ExtractIdPathParameter, func(request uniform.M) uniform.M {
		// todo: use uniform validator to validate fields
		// validator := uniform.NewValidator()
		for key, value := range request {
			switch key {
			default:
				// validator.Error(key, "Unexpected field")
				break
			case "firstName":
				if value == "" {
					// validator.Error("firstName", "May not be empty")
				}
				break
			case "lastName":
				if value == "" {
					// validator.Error("firstName", "May not be empty")
				}
				break
			case "email":
				if value == "" {
					// validator.Error("email", "May not be empty")
				}
				break
			case "mobile":
				if value == "" {
					// validator.Error("mobile", "May not be empty")
				}
				break
			}
		}
		// validator.Check()
		return request
	}, nil)
}