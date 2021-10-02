package actions

/* Example
This is just an example routine for reference
*/

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

func init() {
	_base.Subscribe(_base.TargetAction("ping"), ping)
}

func ping(r uniform.IRequest, p diary.IPage) {
	var model interface{}
	r.Read(model)

	p.Notice("pong", diary.M{
		"timeout": r.Timeout(),
		"remainder": r.Remainder(),
		"hasError": r.HasError(),
		"error": r.Error(),
		"canReply": r.CanReply(),
		"channel": r.Channel(),
		"startedAt": r.StartedAt(),
		"parameters": r.Parameters(),
		"context": r.Context(),
		"model": model,
	})

	if r.CanReply() {
		if err := r.Reply(uniform.Request{}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"err": err,
			})
		}
	}
}