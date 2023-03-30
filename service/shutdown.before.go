package service

import (
	"github.com/go-diary/diary"
	"sync"
)

func ShutdownBefore(shutdown chan bool, group *sync.WaitGroup, p diary.IPage) {
	// todo: add your custom shutdown logic here
}
