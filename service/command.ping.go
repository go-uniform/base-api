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
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("---------- PONG --------------")
	fmt.Println()
	fmt.Println()
	fmt.Println()
}
