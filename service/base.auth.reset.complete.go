package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.reset.complete"), authResetComplete)
}

func authResetComplete(r uniform.IRequest, p diary.IPage) {

}