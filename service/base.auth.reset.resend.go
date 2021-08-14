package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.reset.resend"), authResetResend)
}

func authResetResend(r uniform.IRequest, p diary.IPage) {

}