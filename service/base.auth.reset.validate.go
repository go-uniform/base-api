package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.reset.validate"), authResetValidate)
}

func authResetValidate(r uniform.IRequest, p diary.IPage) {
	if err := p.Scope("auth.reset.validate", func(s diary.IPage) {
		authRequest("reset.validate", r, s)
	}); err != nil {
		panic(err)
	}
}
