package service

import (
	"github.com/go-diary/diary"
	"sync"
)

func ShutdownAfter(shutdown chan bool, group *sync.WaitGroup, p diary.IPage) {
	// todo: add your custom exit logic here
}
