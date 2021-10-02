package commands

/* Example
This is just an example command for reference
*/

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"service/service/_base"
)

func init() {
	_base.Subscribe(_base.TargetCommand("ping"), ping)
}

func ping(r uniform.IRequest, p diary.IPage) {
	context := diary.M{}

	for key, value := range r.Parameters() {
		context[key] = value
	}

	p.Notice("pong", context)

	if r.CanReply() {
		if err := r.Reply(uniform.Request{
			Model: "pong",
		}); err != nil {
			p.Error("reply", err.Error(), diary.M{
				"err": err,
			})
		}
	}
}
