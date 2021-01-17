package service

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	AppClient = "uprate"
	AppProject = "uniform"
)

func Run(page diary.IPage, c uniform.IConn) {
	if err := page.Scope("run", func(p diary.IPage) {
		m := c.Mongo(p, "")

		var model struct {
			Id primitive.ObjectID
			Name string
		}

		m.Insert(time.Second, "uniform", "names", M{
			"name": "Name",
		}, &model, nil)

		fmt.Println(model)
	}); err != nil {
		panic(err)
	}
}