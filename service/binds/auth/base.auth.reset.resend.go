package auth

import (
	"github.com/go-uniform/uniform"
	"net/http"
	"service/service/_base"
)

const AuthResetResend = "auth.reset.resend"

func init() {
	_base.Bind(AuthResetResend, http.MethodPost, "/auth/reset/{id}/resend", _base.ExtractIdPathParameter, func(request uniform.M) uniform.M {
		// todo: use uniform validator to validate fields
		// validator := uniform.NewValidator()
		for key, value := range request {
			switch key {
			default:
				// validator.Error(key, "Unexpected field")
				break
			case "request-id":
				if value == "" {
					// validator.Error("request-id", "May not be empty")
				}
				break
			}
		}
		// validator.Check()
		return request
	}, nil)
}