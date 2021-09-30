package service

import (
	"net/http"
)

const AuthResetResend = "auth.reset.resend"

func init() {
	bind(AuthResetResend, http.MethodPost, "/auth/reset/{id}/resend", extractIdPathParameter, func(request M) M {
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