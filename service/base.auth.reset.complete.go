package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.reset.complete"), authResetComplete)
}

func authResetComplete(r uniform.IRequest, p diary.IPage) {
	if err := p.Scope("auth.reset.complete", func(s diary.IPage) {
		authRequest("reset.complete", r, s)
	}); err != nil {
		panic(err)
	}
}