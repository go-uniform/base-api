package auth

import (
	"github.com/go-uniform/uniform"
	"net/http"
	"service/service/_base"
)

func init() {
	_base.Bind("auth.login", http.MethodPost, "/auth/login", nil, func(request uniform.M) uniform.M {
		fields := []string{
			"type",
			"identifier",
			"password",
		}
		requiredFields := []string{
			"type",
			"identifier",
			"password",
		}

		var validator uniform.IValidator
		validator, request = uniform.RequestValidator(request, fields...)
		request = validator.Required(request, requiredFields...)
		validator.Validate()

		return request
	}, nil)
}
