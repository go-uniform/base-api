package service

import (
	"net/http"
)

const AuthResetValidate = "auth.reset.validate"

func init() {
	bind(AuthResetValidate, http.MethodPost, "/auth/reset/{id}/validate", extractIdPathParameter, func(request M) M {
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