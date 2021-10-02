package auth

import (
	"github.com/go-uniform/uniform"
	"net/http"
	"service/service/_base"
)

const TopicAuthLoginCodeValidate = "auth.login.code-validate"

func init() {
	_base.Bind(TopicAuthLoginCodeValidate, http.MethodPost, "/auth/login/code-validate", nil, func(request uniform.M) uniform.M {
		// todo: use uniform validator to validate fields
		// validator := uniform.NewValidator()
		for key, value := range request {
			switch key {
			default:
				// validator.Error(key, "Unexpected field")
				break
			case "token":
				if value == "" {
					// validator.Error("token", "May not be empty")
				}
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