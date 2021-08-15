package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local("action.administrators.list"), administratorsList)
}

func administratorsList(r uniform.IRequest, p diary.IPage) {
}