package service

import (
	"net/http"
)

const TopicAuthReset = "auth.reset"

func init() {
	bind(TopicAuthReset, http.MethodPost, "/auth/reset", nil, func(request M) M {
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
			case "method":
				if value == "" {
					// validator.Error("method", "May not be empty")
				}
				break
			case "channel":
				if value == "" {
					// validator.Error("channel", "May not be empty")
				}
				break
			}
		}
		// validator.Check()
		return request
	}, nil)
}