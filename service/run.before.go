package service

import (
	"github.com/go-diary/diary"
	"sync"
)

func RunBefore(shutdown chan bool, group *sync.WaitGroup, p diary.IPage) {
}