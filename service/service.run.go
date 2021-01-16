package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"time"
)

const (
	AppClient = "uprate"
	AppProject = "uniform"
)

func Run(page diary.IPage, c uniform.IConn) {
	if err := page.Scope("run", func(p diary.IPage) {
		m := c.Mongo(p, "")

		m.Insert(time.Second, "pay-curve", "names", M{
			"name": "Name",
		}, nil, nil)
	}); err != nil {
		panic(err)
	}
}