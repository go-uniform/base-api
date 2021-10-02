package auth

import (
	"net/http"
	"service/service/_base"
)

const TopicAuthLoginCodeResend = "auth.login.code-resend"

func init() {
	_base.bind(TopicAuthLoginCodeResend, http.MethodPost, "/auth/login/code-resend", nil, func(request _base.M) _base.M {
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
			}
		}
		// validator.Check()
		return request
	}, nil)
}