package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.auth.reset.validate"), authResetValidate)
}

func authResetValidate(r uniform.IRequest, p diary.IPage) {

}
