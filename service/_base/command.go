package _base

import (
	"fmt"
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
	"os"
	"runtime/debug"
	"service/service/info"
	"time"
)

func Command(cmd string, timeout time.Duration, natsUri string, natsOptions []nats.Option, args uniform.P, responseHandler func([]byte)) {
	// no diary for commands since we want response data to be only thing is stdout
	InitializeNoDiary()

	natsConn, err := nats.Connect(natsUri, natsOptions...)
	if err != nil {
		panic(err)
	}
	info.Conn, err = uniform.ConnectorNats(info.Diary, natsConn)
	if err != nil {
		panic(err)
	}

	defer info.Conn.Close()

	info.Diary.Page(-1, info.TraceRate, true, info.AppName, nil, "", "", nil, func(p diary.IPage) {
		if err := info.Conn.Request(p, TargetCommand(cmd), timeout, uniform.Request{
			Model: args,
		}, func(r uniform.IRequest, p diary.IPage) {
			if r.HasError() {
				fmt.Println(r.Error())
				os.Exit(1)
			}
			if responseHandler != nil {
				var data []byte
				r.Read(&data)
				responseHandler(data)
			}
		}); err != nil {
			fmt.Println(err.Error())
			fmt.Println(string(debug.Stack()))
			os.Exit(1)
		}
	})

	_ = info.Conn.Drain()
}
