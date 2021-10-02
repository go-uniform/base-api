package auth

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service"
	"strings"
)

func authRequest(action string, r uniform.IRequest, p diary.IPage) {
	var response service.M
	if err := r.Conn().Request(p, fmt.Sprintf("auth.%s", strings.TrimPrefix(action, "auth.")), r.Remainder(), uniform.Request{
		Model: service.M{
			"group": AppProject,
		},
	}, func(r uniform.IRequest, p diary.IPage) {
		if r.HasError() {
			panic(r.Error())
		}
		r.Read(response)
	}); err != nil {
		panic(err)
	}

	if err := r.Reply(uniform.Request{
		Model: response,
	}); err != nil {
		p.Error("reply", err.Error(), diary.M{
			"error": err,
		})
	}
}