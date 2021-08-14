package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(routine("println"), println)
}

func println(r uniform.IRequest, p diary.IPage) {
	var message string
	r.Read(&message)

	fmt.Println("-----------------------------------------------------------")
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println(message)
	fmt.Println()
	fmt.Println()
	fmt.Println()
	fmt.Println("-----------------------------------------------------------")

	if r.CanReply() {
		if err := r.Reply(uniform.Request{}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"err": err,
			})
		}
	}
}