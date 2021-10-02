package auth

import (
	"net/http"
	"service/service"
)

const TopicAuthLoginCodeValidate = "auth.login.code-validate"

func init() {
	service.bind(TopicAuthLoginCodeValidate, http.MethodPost, "/auth/login/code-validate", nil, func(request service.M) service.M {
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