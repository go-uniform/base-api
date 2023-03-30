package _base

import (
	"fmt"
	"github.com/go-diary/diary"
	"service/service/info"
	"strings"
)

func InitializeNoDiary() {
	info.Diary = diary.Dear(info.AppClient, info.AppProject, info.AppName, nil, info.AppRepository, info.AppCommit, []string{info.AppVersion}, nil, diary.LevelTrace, func(log diary.Log) {
		// future: we may still want to write these logs to the remote log server

		// don't write logs to stdout
		return
	})
}

func InitializeDiary(test, virtual bool, level string, rate int) {
	handler := diary.HumanReadableHandler
	if test {
		// test mode is used for creating an automated testing environment
		fmt.Println("entering test mode")
		handler = nil
	}

	lvl := diary.ConvertFromTextLevel(level)
	if !diary.IsValidLevel(lvl) {
		panic(fmt.Sprintf("level must be one of the following values: %s", strings.Join(diary.TextLevels, ", ")))
	}
	info.TraceRate = rate

	info.Diary = diary.Dear(info.AppClient, info.AppProject, info.AppName, nil, info.AppRepository, info.AppCommit, []string{info.AppVersion}, nil, lvl, handler)
}
