package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.login.code-resend"), authLoginCodeResend)
}

func authLoginCodeResend(r uniform.IRequest, p diary.IPage) {
	if err := p.Scope("auth.login.code-resend", func(s diary.IPage) {
		authRequest("login.code-resend", r, s)
	}); err != nil {
		panic(err)
	}
}