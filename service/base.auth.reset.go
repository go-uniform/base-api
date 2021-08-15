package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.reset"), authReset)
}

func authReset(r uniform.IRequest, p diary.IPage) {
	if err := p.Scope("auth.reset", func(s diary.IPage) {
		authRequest("reset", r, s)
	}); err != nil {
		panic(err)
	}
}