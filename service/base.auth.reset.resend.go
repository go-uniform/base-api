package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.reset.resend"), authResetResend)
}

func authResetResend(r uniform.IRequest, p diary.IPage) {
	if err := p.Scope("auth.reset.resend", func(s diary.IPage) {
		authRequest("reset.resend", r, s)
	}); err != nil {
		panic(err)
	}
}