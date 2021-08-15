package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.login.code-validate"), authLoginCodeValidate)
}

func authLoginCodeValidate(r uniform.IRequest, p diary.IPage) {
	if err := p.Scope("auth.login.code-validate", func(s diary.IPage) {
		authRequest("login.code-validate", r, s)
	}); err != nil {
		panic(err)
	}
}