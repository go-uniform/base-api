package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.login"), authLogin)
}

func authLogin(r uniform.IRequest, p diary.IPage) {

}