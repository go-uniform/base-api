package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(command("ping"), ping)
}

func ping(r uniform.IRequest, p diary.IPage) {
	fmt.Println("\n\n\n---------- PONG --------------\n\n\n")
}
