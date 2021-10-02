package _base

import (
	"github.com/go-diary/diary"
	"github.com/go-uniform/uniform"
	"github.com/nats-io/go-nats"
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
	c, err = uniform.ConnectorNats(d, natsConn)
	if err != nil {
		panic(err)
	}

	defer c.Close()

	d.Page(-1, traceRate, true, info.AppName, nil, "", "", nil, func(p diary.IPage) {
		if err := c.Request(p, TargetCommand(cmd), timeout, uniform.Request{
			Parameters: args,
		}, func(r uniform.IRequest, p diary.IPage) {
			if responseHandler != nil {
				var data []byte
				r.Read(&data)
				responseHandler(data)
			}
		}); err != nil {
			panic(err)
		}
	})

	_ = c.Drain()
}
