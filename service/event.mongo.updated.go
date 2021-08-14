package service

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
)

func init() {
	subscribe(event("mongo", "updated"), eventMongoUpdated)
}

func eventMongoUpdated(r uniform.IRequest, p diary.IPage) {
	// todo: react based on which database and collection has been updated
}