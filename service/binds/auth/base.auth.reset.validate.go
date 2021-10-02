package auth

import (
	"github.com/go-uniform/uniform"
	"net/http"
	"service/service/_base"
)

const AuthResetValidate = "auth.reset.validate"

func init() {
	_base.Bind(AuthResetValidate, http.MethodPost, "/auth/reset/{id}/validate", _base.ExtractIdPathParameter, func(request uniform.M) uniform.M {
		// todo: use uniform validator to validate fields
		// validator := uniform.NewValidator()
		for key, value := range request {
			switch key {
			default:
				// validator.Error(key, "Unexpected field")
				break
			case "code":
				if value == "" {
					// validator.Error("code", "May not be empty")
				}
				break
			}
		}
		// validator.Check()
		return request
	}, nil)
}