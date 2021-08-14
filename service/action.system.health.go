package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(local(system("health")), health)
}

func health(r uniform.IRequest, p diary.IPage) {
	// todo: check service health and throw error if unhealthy

	if r.CanReply() {
		if err := r.Reply(uniform.Request{}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"err": err,
			})
		}
	}
}
