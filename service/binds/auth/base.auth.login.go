package auth

import (
	"net/http"
	"service/service"
)

const TopicAuthLogin = "auth.login"

func init() {
	service.bind(TopicAuthLogin, http.MethodPost, "/auth/login", nil, func(request service.M) service.M {
		// todo: use uniform validator to validate fields
		// validator := uniform.NewValidator()
		for key, value := range request {
			switch key {
			default:
				// validator.Error(key, "Unexpected field")
				break
			case "type":
				if value == "" {
					// validator.Error("type", "May not be empty")
				}
				break
			case "identifier":
				if value == "" {
					// validator.Error("identifier", "May not be empty")
				}
				break
			case "password":
				if value == "" {
					// validator.Error("password", "May not be empty")
				}
				break
			}
		}
		// validator.Check()
		return request
	}, nil)
}