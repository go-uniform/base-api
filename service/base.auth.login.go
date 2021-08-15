package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.login"), authLogin)
}

func authLogin(r uniform.IRequest, p diary.IPage) {
	if err := p.Scope("auth.login", func(s diary.IPage) {
		authRequest("login", r, s)
	}); err != nil {
		panic(err)
	}
}